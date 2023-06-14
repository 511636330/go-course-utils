package course

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
	"gitlab.qkids.com/group-api-common/go-redis.git"
)

var storage *cache.Cache

const (
	connection         string = "course"
	QkdisBasicSeriesId int8   = 1
	AISeriesId         int8   = 22
	ArtSeriesId        int8   = 17
	QluckSeriesId      int8   = 8
)

func init() {
	storage = cache.New(5*time.Minute, 10*time.Minute)
}

func GetCourses() (courses []CourseBrief, err error) {
	cacheKey := "COURSES"
	if c, found := storage.Get(cacheKey); found {
		courses = *(c.(*[]CourseBrief))
		return
	}
	str, err := GetCacheFromRedis(cacheKey)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(str), &courses)

	if err == nil {
		storage.Set(cacheKey, &courses, GetCacheDuration())
		return
	}

	return
}

func GetCoursesByLevelAndSeries(level, seriesId int8) (wantCourses []CourseBrief) {
	courses, err := GetCourses()
	if err != nil {
		return
	}

	for _, course := range courses {
		if course.Level != level {
			continue
		}
		if seriesId > 0 && course.SeriesId != seriesId {
			continue
		}
		wantCourses = append(wantCourses, course)
	}

	return
}

func GetCourseBrief(id int) (course CourseBrief) {
	courses, err := GetCourses()
	if err != nil {
		return
	}
	for _, c := range courses {
		if c.Id == id {
			return c
		}
	}
	return
}

func GetSeriesIds(courseIds []int) (ids []int8) {
	if len(courseIds) == 0 {
		return
	}
	courses, err := GetCourses()
	if err != nil {
		return
	}

	courseIdsMap := map[int]byte{}
	idsMap := map[int8]byte{}

	for _, courseId := range courseIds {
		courseIdsMap[courseId] = 0
	}
	for _, course := range courses {
		if _, ok := courseIdsMap[course.Id]; !ok {
			continue
		}
		idsMap[course.SeriesId] = 0
		if len(idsMap) != len(ids) {
			ids = append(ids, course.SeriesId)
		}
	}

	return
}

func GetCourse(id int) (course Course, err error) {
	cacheKey := fmt.Sprintf("COURSE:%d", id)
	if c, found := storage.Get(cacheKey); found {
		course = *(c.(*Course))
		return
	}

	str, err := GetCacheFromRedis(cacheKey)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(str), &course)

	if err == nil {
		storage.Set(cacheKey, &course, GetCacheDuration())
	}

	return
}

func (course Course) IsContainLesson(lessonId int) bool {
	for _, chapter := range course.Chapters {
		for _, lesson := range chapter.Lessons {
			if lesson.Id == lessonId {
				return true
			}
		}
	}

	return false
}

func (course Course) GetLesson(id int) (lesson ChapterLesson, found bool) {
	for _, chapter := range course.Chapters {
		for _, lesson := range chapter.Lessons {
			if lesson.Id == id {
				return lesson, true
			}
		}
	}
	return
}

func (course Course) GetLessons(ids []int) (lessons []ChapterLesson) {
	for _, id := range ids {
		if lesson, found := course.GetLesson(id); found {
			lessons = append(lessons, lesson)
		}
	}
	return
}

func (course Course) GetLessonIdsInOrder() (ids []int) {
	for _, chapter := range course.Chapters {
		for _, lesson := range chapter.Lessons {
			ids = append(ids, lesson.Id)
		}
	}
	return
}

func (course Course) GetChapterLessonIdsInOrder(chapterId int) (ids []int) {
	for _, chapter := range course.Chapters {
		if chapter.Id != chapterId {
			continue
		}
		for _, lesson := range chapter.Lessons {
			ids = append(ids, lesson.Id)
		}
	}
	return
}

func (course Course) IsUpgradeTestPublished() bool {
	return course.UpgradeTest
}

func (course Course) GetChapterReviewLessonId(chapterId int) int {
	for _, chapter := range course.Chapters {
		if chapter.Id != chapterId {
			continue
		}
		for _, lesson := range chapter.Lessons {
			if lesson.Review {
				return lesson.Id
			}
		}
	}
	return 0
}

func GetLesson(id int) (lesson Lesson, err error) {
	cacheKey := fmt.Sprintf("LESSON:%d", id)
	if l, found := storage.Get(cacheKey); found {
		lesson = *(l.(*Lesson))
		return
	}
	str, err := GetCacheFromRedis(cacheKey)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(str), &lesson)
	if err == nil {
		storage.Set(cacheKey, &lesson, GetCacheDuration())
	}

	return
}

func GetCourseIdFromLesson(lessonId int) int {
	cacheKey := fmt.Sprintf("LESSON:%d:COURSE_ID", lessonId)
	if courseId, found := storage.Get(cacheKey); found {
		return courseId.(int)
	}

	courseIdString, err := GetCacheFromRedis(cacheKey)
	if err != nil {
		courseId := cast.ToInt(courseIdString)

		if courseId > 0 {
			storage.Set(cacheKey, courseId, GetCacheDuration())
		}
		return courseId
	}

	return 0
}

func GetCourseTests(courseId int) []Test {
	var courseTest CourseTest
	cacheKey := fmt.Sprintf("COURSE:%d:TESTS", courseId)
	if ts, found := storage.Get(cacheKey); found {
		courseTest = *(ts.(*CourseTest))
		return courseTest.Tests
	}

	str, err := GetCacheFromRedis(cacheKey)

	if err != nil {
		return []Test{}
	}

	err = json.Unmarshal([]byte(str), &courseTest)

	if err == nil {
		storage.Set(cacheKey, &courseTest, GetCacheDuration())
	}

	return courseTest.Tests
}

func GetLessonTests(lessonId int) []Test {
	var lessonTest LessonTest
	cacheKey := fmt.Sprintf("LESSON:%d:TESTS", lessonId)
	if ts, found := storage.Get(cacheKey); found {
		lessonTest = *(ts.(*LessonTest))
		return lessonTest.Tests
	}

	str, err := GetCacheFromRedis(cacheKey)

	if err != nil {
		return []Test{}
	}

	err = json.Unmarshal([]byte(str), &lessonTest)

	if err == nil {
		storage.Set(cacheKey, &lessonTest, GetCacheDuration())
	}

	return lessonTest.Tests
}

func GetCourseIdFromSeriesId(seriesId int8) (courseIds []int) {
	if courses, err := GetCourses(); err == nil {
		for _, course := range courses {
			if course.SeriesId == seriesId {
				courseIds = append(courseIds, course.Id)
			}
		}
	}
	return
}

func GetAppBasicSeriesId(appId int8) int8 {
	switch appId {
	case 0:
		return QkdisBasicSeriesId
	case 1:
		return AISeriesId
	case 2:
		return ArtSeriesId
	default:
		return QkdisBasicSeriesId
	}
}

func GetCacheFromRedis(cacheKey string) (str string, err error) {
	redisClient := redis.GetClient(connection)
	str, err = redisClient.Get(context.TODO(), cacheKey).Result()
	return

}

func GetCacheDuration() time.Duration {
	hour := time.Now().Hour()

	if (hour >= 8 && hour <= 13) || (hour >= 17 && hour <= 21) || hour >= 23 || hour <= 5 {
		return time.Hour
	}
	return 5 * time.Minute
}

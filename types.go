package course

type CourseBrief struct {
	Id            int
	Name          string
	NameEn        string
	NameTr        string
	Level         int8
	Published     bool
	Abbreviation  string
	SeriesId      int8
	SeriesName    string
	UpgradeTest   bool
	CategoryId    int8
	LogoUri       string
	AppLogoUri    string
	Description   string
	DescriptionEn string
}

type ChapterLesson struct {
	Id            int
	Name          string
	Published     bool
	Trial         bool
	Review        bool
	LogoUri       string
	LogoHDUri     string
	ThemeColor    string
	Description   string
	DescriptionCN string
}

type CourseChapter struct {
	Id      int
	Name    string
	Lessons []ChapterLesson
}

type Course struct {
	Id            int
	Name          string
	NameEn        string
	NameTr        string
	Level         int8
	Published     bool
	Abbreviation  string
	SeriesId      int8
	SeriesName    string
	UpgradeTest   bool
	CategoryId    int8
	LogoUri       string
	AppLogoUri    string
	Description   string
	DescriptionEn string
	Chapters      []CourseChapter
}

type ChapterCourse struct {
	Id           int
	Name         string
	NameEn       string
	NameTr       string
	Level        int8
	Abbreviation string
	SeriesId     int8
	SeriesName   string
}

type LessonChapter struct {
	Id     int
	Name   string
	Course ChapterCourse
}

type Lesson struct {
	Id         int
	Name       string
	Published  bool
	Trial      bool
	Review     bool
	LogoUri    string
	LogoHDUri  string
	ThemeColor string
	Chapter    LessonChapter
}

type SubTest struct {
	Id int
}

type Test struct {
	Id       int
	Category string
	Tests    []SubTest
}

type CourseTest struct {
	Id    int
	Tests []Test
}

type LessonTest struct {
	Id    int
	Tests []Test
}

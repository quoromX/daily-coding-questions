package database

import "time"

type User struct {
	ID           string    `json:"id"`
	Handle       string    `json:"handle"`
	DisplayName  string    `json:"displayName"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Rank         string    `json:"rank"`
	XP           int       `json:"xp"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Problem struct {
	ID             string        `json:"id"`
	Slug           string        `json:"slug"`
	Title          string        `json:"title"`
	Difficulty     string        `json:"difficulty"`
	Gate           string        `json:"gate"`
	Summary        string        `json:"summary"`
	Body           string        `json:"body"`
	Constraints    []string      `json:"constraints"`
	Examples       []Example     `json:"examples"`
	Tags           []string      `json:"tags"`
	StarterCode    []StarterCode `json:"starterCode"`
	TestCases      []TestCase    `json:"-"`
	AcceptanceRate int           `json:"acceptanceRate"`
	XP             int           `json:"xp"`
}

type Example struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Note   string `json:"note"`
}

type StarterCode struct {
	LanguageID int    `json:"languageId"`
	Language   string `json:"language"`
	Code       string `json:"code"`
}

type TestCase struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expectedOutput"`
	Hidden         bool   `json:"hidden"`
}

type Submission struct {
	ID          string             `json:"id"`
	UserID      string             `json:"userId"`
	ProblemSlug string             `json:"problemSlug"`
	ProblemName string             `json:"problemName"`
	LanguageID  int                `json:"languageId"`
	Language    string             `json:"language"`
	SourceCode  string             `json:"sourceCode"`
	Status      string             `json:"status"`
	Score       int                `json:"score"`
	RuntimeMS   int                `json:"runtimeMs"`
	MemoryKB    int                `json:"memoryKb"`
	Results     []SubmissionResult `json:"results"`
	CreatedAt   time.Time          `json:"createdAt"`
}

type SubmissionResult struct {
	TestName      string `json:"testName"`
	Status        string `json:"status"`
	Stdout        string `json:"stdout"`
	Stderr        string `json:"stderr"`
	CompileOutput string `json:"compileOutput"`
	RuntimeMS     int    `json:"runtimeMs"`
	MemoryKB      int    `json:"memoryKb"`
}

type Dashboard struct {
	User              User         `json:"user"`
	Solved            int          `json:"solved"`
	Attempted         int          `json:"attempted"`
	Streak            int          `json:"streak"`
	Accuracy          int          `json:"accuracy"`
	BestLanguage      string       `json:"bestLanguage"`
	DifficultySpread  []Metric     `json:"difficultySpread"`
	RecentSubmissions []Submission `json:"recentSubmissions"`
	Recommended       []Problem    `json:"recommended"`
}

type Metric struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

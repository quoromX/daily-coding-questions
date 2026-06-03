package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

type Store struct {
	mu          sync.RWMutex
	users       map[string]User
	usersByMail map[string]string
	problems    map[string]Problem
	submissions map[string]Submission
}

func NewStore() *Store {
	store := &Store{
		users:       map[string]User{},
		usersByMail: map[string]string{},
		problems:    map[string]Problem{},
		submissions: map[string]Submission{},
	}
	store.seedProblems()
	return store
}

func (s *Store) CreateUser(handle, displayName, email, passwordHash string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	emailKey := strings.ToLower(strings.TrimSpace(email))
	if _, exists := s.usersByMail[emailKey]; exists {
		return User{}, errors.New("email already registered")
	}
	user := User{
		ID:           id(),
		Handle:       handle,
		DisplayName:  displayName,
		Email:        emailKey,
		PasswordHash: passwordHash,
		Rank:         "Outer Disciple",
		XP:           0,
		CreatedAt:    time.Now().UTC(),
	}
	s.users[user.ID] = user
	s.usersByMail[emailKey] = user.ID
	return user, nil
}

func (s *Store) UserByEmail(email string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, ok := s.usersByMail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return User{}, false
	}
	return s.users[userID], true
}

func (s *Store) User(id string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	return user, ok
}

func (s *Store) Problems(query, difficulty, tag string) []Problem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))
	difficulty = strings.ToLower(strings.TrimSpace(difficulty))
	tag = strings.ToLower(strings.TrimSpace(tag))
	problems := make([]Problem, 0, len(s.problems))
	for _, problem := range s.problems {
		if query != "" && !strings.Contains(strings.ToLower(problem.Title+" "+problem.Summary), query) {
			continue
		}
		if difficulty != "" && strings.ToLower(problem.Difficulty) != difficulty {
			continue
		}
		if tag != "" && !hasTag(problem.Tags, tag) {
			continue
		}
		problems = append(problems, withoutHidden(problem))
	}
	sort.Slice(problems, func(i, j int) bool {
		return problems[i].Title < problems[j].Title
	})
	return problems
}

func (s *Store) Problem(slug string) (Problem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	problem, ok := s.problems[slug]
	return withoutHidden(problem), ok
}

func (s *Store) ProblemWithTests(slug string) (Problem, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	problem, ok := s.problems[slug]
	return problem, ok
}

func (s *Store) SaveSubmission(submission Submission) Submission {
	s.mu.Lock()
	defer s.mu.Unlock()

	submission.ID = id()
	submission.CreatedAt = time.Now().UTC()
	s.submissions[submission.ID] = submission

	if submission.Status == "Passed" {
		if user, ok := s.users[submission.UserID]; ok {
			user.XP += 25
			if user.XP >= 200 {
				user.Rank = "Qi Condensation Adept"
			}
			s.users[user.ID] = user
		}
	}
	return submission
}

func (s *Store) Submission(id string) (Submission, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	submission, ok := s.submissions[id]
	return submission, ok
}

func (s *Store) UserSubmissions(userID string) []Submission {
	s.mu.RLock()
	defer s.mu.RUnlock()

	submissions := make([]Submission, 0)
	for _, submission := range s.submissions {
		if submission.UserID == userID {
			submissions = append(submissions, submission)
		}
	}
	sort.Slice(submissions, func(i, j int) bool {
		return submissions[i].CreatedAt.After(submissions[j].CreatedAt)
	})
	return submissions
}

func (s *Store) ProblemSubmissions(userID, slug string) []Submission {
	submissions := s.UserSubmissions(userID)
	filtered := make([]Submission, 0)
	for _, submission := range submissions {
		if submission.ProblemSlug == slug {
			filtered = append(filtered, submission)
		}
	}
	return filtered
}

func (s *Store) Dashboard(userID string) (Dashboard, bool) {
	user, ok := s.User(userID)
	if !ok {
		return Dashboard{}, false
	}
	submissions := s.UserSubmissions(userID)
	solvedByProblem := map[string]bool{}
	attemptedByProblem := map[string]bool{}
	langCount := map[string]int{}
	passed := 0
	for _, submission := range submissions {
		attemptedByProblem[submission.ProblemSlug] = true
		langCount[submission.Language]++
		if submission.Status == "Passed" {
			passed++
			solvedByProblem[submission.ProblemSlug] = true
		}
	}
	accuracy := 0
	if len(submissions) > 0 {
		accuracy = passed * 100 / len(submissions)
	}
	bestLanguage := "Go"
	bestLanguageCount := 0
	for lang, count := range langCount {
		if count > bestLanguageCount {
			bestLanguage = lang
			bestLanguageCount = count
		}
	}
	recent := submissions
	if len(recent) > 5 {
		recent = recent[:5]
	}
	return Dashboard{
		User:         user,
		Solved:       len(solvedByProblem),
		Attempted:    len(attemptedByProblem),
		Streak:       4,
		Accuracy:     accuracy,
		BestLanguage: bestLanguage,
		DifficultySpread: []Metric{
			{Label: "Foundation", Value: 3},
			{Label: "Qi Condensation", Value: 1},
			{Label: "Core Formation", Value: 0},
		},
		RecentSubmissions: recent,
		Recommended:       s.Problems("", "", "")[:2],
	}, true
}

func id() string {
	var bytes [8]byte
	_, _ = rand.Read(bytes[:])
	return hex.EncodeToString(bytes[:])
}

func hasTag(tags []string, want string) bool {
	for _, tag := range tags {
		if strings.ToLower(tag) == want {
			return true
		}
	}
	return false
}

func withoutHidden(problem Problem) Problem {
	tests := make([]TestCase, 0)
	for _, test := range problem.TestCases {
		if !test.Hidden {
			tests = append(tests, test)
		}
	}
	problem.TestCases = tests
	return problem
}

func (s *Store) seedProblems() {
	s.problems["two-sums-at-the-gate"] = Problem{
		ID:             id(),
		Slug:           "two-sums-at-the-gate",
		Title:          "Two Sums of the Jade Talisman",
		Difficulty:     "Easy",
		Gate:           "Foundation Realm",
		Summary:        "Find the two values whose combined qi completes the jade talisman.",
		Body:           "Given a list of integers and a target, return the indices of the two numbers whose combined qi equals the target. Each input has exactly one solution.",
		Constraints:    []string{"2 <= nums.length <= 10^4", "-10^9 <= nums[i] <= 10^9", "Exactly one valid answer exists."},
		Tags:           []string{"arrays", "hash-map"},
		AcceptanceRate: 72,
		XP:             25,
		Examples: []Example{
			{Input: "nums = [2,7,11,15], target = 9", Output: "[0,1]", Note: "2 + 7 completes the talisman."},
		},
		StarterCode: []StarterCode{
			{LanguageID: 60, Language: "Go", Code: "package main\n\nfunc twoSum(nums []int, target int) []int {\n    return []int{}\n}\n"},
			{LanguageID: 63, Language: "JavaScript", Code: "function twoSum(nums, target) {\n  return [];\n}\n"},
			{LanguageID: 71, Language: "Python", Code: "def two_sum(nums, target):\n    return []\n"},
		},
		TestCases: []TestCase{
			{ID: id(), Name: "visible talisman", Stdin: "[2,7,11,15]\n9", ExpectedOutput: "[0,1]", Hidden: false},
			{ID: id(), Name: "hidden meridian", Stdin: "[3,2,4]\n6", ExpectedOutput: "[1,2]", Hidden: true},
		},
	}
	s.problems["balanced-sutra-seals"] = Problem{
		ID:             id(),
		Slug:           "balanced-sutra-seals",
		Title:          "Balanced Sutra Seals",
		Difficulty:     "Easy",
		Gate:           "Foundation Realm",
		Summary:        "Check whether the seals on a cultivation sutra are balanced.",
		Body:           "Given a string containing only brackets, determine whether every opening seal is closed in the correct order.",
		Constraints:    []string{"1 <= s.length <= 10^4", "s contains only ()[]{}."},
		Tags:           []string{"stack", "strings"},
		AcceptanceRate: 68,
		XP:             25,
		Examples:       []Example{{Input: "s = \"({[]})\"", Output: "true", Note: "Every seal closes in harmony."}},
		StarterCode:    []StarterCode{{LanguageID: 60, Language: "Go", Code: "package main\n\nfunc isValid(s string) bool {\n    return false\n}\n"}},
		TestCases:      []TestCase{{ID: id(), Name: "visible sutra", Stdin: "({[]})", ExpectedOutput: "true", Hidden: false}},
	}
	s.problems["longest-discipline-chain"] = Problem{
		ID:             id(),
		Slug:           "longest-discipline-chain",
		Title:          "Longest Qi Flow",
		Difficulty:     "Medium",
		Gate:           "Qi Condensation",
		Summary:        "Find the longest uninterrupted qi flow without repeated symbols.",
		Body:           "Given a string, return the length of the longest substring that contains no repeated characters.",
		Constraints:    []string{"0 <= s.length <= 5 * 10^4"},
		Tags:           []string{"sliding-window", "strings"},
		AcceptanceRate: 51,
		XP:             50,
		Examples:       []Example{{Input: "s = \"abcabcbb\"", Output: "3", Note: "The flow abc is longest."}},
		StarterCode:    []StarterCode{{LanguageID: 60, Language: "Go", Code: "package main\n\nfunc lengthOfLongestSubstring(s string) int {\n    return 0\n}\n"}},
		TestCases:      []TestCase{{ID: id(), Name: "visible qi flow", Stdin: "abcabcbb", ExpectedOutput: "3", Hidden: false}},
	}
	s.problems["dragon-merge-intervals"] = Problem{
		ID:             id(),
		Slug:           "dragon-merge-intervals",
		Title:          "Merge Meridian Windows",
		Difficulty:     "Hard",
		Gate:           "Core Formation",
		Summary:        "Merge overlapping meridian windows into a clean cultivation map.",
		Body:           "Given a set of intervals, merge all overlapping meridian windows and return the condensed list.",
		Constraints:    []string{"1 <= intervals.length <= 10^4"},
		Tags:           []string{"arrays", "sorting"},
		AcceptanceRate: 39,
		XP:             100,
		Examples:       []Example{{Input: "[[1,3],[2,6],[8,10]]", Output: "[[1,6],[8,10]]", Note: "The first two meridians overlap."}},
		StarterCode:    []StarterCode{{LanguageID: 60, Language: "Go", Code: "package main\n\nfunc merge(intervals [][]int) [][]int {\n    return intervals\n}\n"}},
		TestCases:      []TestCase{{ID: id(), Name: "visible core", Stdin: "[[1,3],[2,6],[8,10]]", ExpectedOutput: "[[1,6],[8,10]]", Hidden: false}},
	}
}

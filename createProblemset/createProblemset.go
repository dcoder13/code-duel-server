package createProblemset

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type Problem struct {
	ContestId	string
	ProblemId	string
	Name		string
	Statement	string
	Constraints	string
	Rating		string
}

type Submission struct {
	Problem Problem `json:"problem"`
	Verdict string  `json:"verdict"`
}

type APIResponse[T any] struct {
	Status string `json:"status"`
	Result T      `json:"result"`
}

func problemKey(p Problem) string {
	return fmt.Sprintf("%s-%s", p.ContestId, p.ProblemId)
}

func fetchUserSolvedProblems(handle string) (map[string]bool, error) {
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s", handle)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result APIResponse[[]Submission]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	solved := make(map[string]bool)
	for _, sub := range result.Result {
		if sub.Verdict == "OK" {
			solved[problemKey(sub.Problem)] = true
		}
	}
	return solved, nil
}

func fetchAllProblems() ([]Problem, error) {
	url := "https://codeforces.com/api/problemset.problems"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw struct {
		Status string `json:"status"`
		Result struct {
			Problems []Problem `json:"problems"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return raw.Result.Problems, nil
}

func bucketProblems(problems []Problem, solved1, solved2 map[string]bool) (easy, medium, hard []Problem) {
	for _, p := range problems {
		key := problemKey(p)
		if solved1[key] || solved2[key] || p.Rating == "0" {
			continue
		}

		switch {
		case p.Rating <= "1200":
			easy = append(easy, p)
		case p.Rating <= "1700":
			medium = append(medium, p)
		default:
			hard = append(hard, p)
		}
	}
	return
}

func createProblemset(handle1 string, handle2 string) (problem [][]Problem) {
	solved1, err := fetchUserSolvedProblems(handle1)
	if err != nil {
		panic("error fetching user 1: " + err.Error())
	}

	solved2, err := fetchUserSolvedProblems(handle2)
	if err != nil {
		panic("error fetching user 2: " + err.Error())
	}

	allProblems, err := fetchAllProblems()
	if err != nil {
		panic("error fetching problems: " + err.Error())
	}

	easy, medium, hard := bucketProblems(allProblems, solved1, solved2)
	problem = [][]Problem{easy, medium, hard}

	return problem
}
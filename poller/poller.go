package poller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dcoder13/code-duel-server/createProblemset"
)

type Submission struct {
    ID       int     `json:"id"`
    Problem  createProblemset.Problem `json:"problem"`
    Verdict  string  `json:"verdict"`
}

func pollVerdict(handle string, contestId string, ProblemId string, maxAttempts int, interval time.Duration) string {
    api := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s", handle)

    for i := 0; i < maxAttempts; i++ {
        resp, err := http.Get(api)
        if err != nil {
            return fmt.Sprintln("API error:", err) 
        }
        body, _ := io.ReadAll(resp.Body)
        resp.Body.Close()

        var result struct {
            Status string       `json:"status"`
            Result []Submission `json:"result"`
        }
        json.Unmarshal(body, &result)

        for _, sub := range result.Result {
            if sub.Problem.ContestId == contestId && sub.Problem.ProblemId == ProblemId {
                if sub.Verdict != "" {
                    fmt.Println("Verdict found:", sub.Verdict)
                    return sub.Verdict
                }
            }
        }

        fmt.Println("Waiting for verdict...")
        time.Sleep(interval)
    }

    fmt.Println("Timeout: Verdict not found")
	return "NIL"
}



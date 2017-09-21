package engine

import (
	"strconv"
	"strings"

	"github.com/imyoudu/goreporter/linters/deadcode"
)

type StrategyDeadCode struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyDeadCode) GetName() string {
	return "Deadcode"
}

func (s *StrategyDeadCode) GetDescription() string {
	return "All useless code, or never obsolete obsolete code."
}

func (s *StrategyDeadCode) GetWeight() float64 {
	return 0.04
}

// linterDead provides a function that will scans all useless code, or never
// obsolete obsolete code.It will extract from the linter need to convert
// the data.The result will be saved in the r's attributes.
func (s *StrategyDeadCode) Compute(parameters StrategyParameter) (summaries Summaries) {
	summaries = NewSummaries()

	deadcodes := deadcode.DeadCode(parameters.ProjectPath)
	sumProcessNumber := int64(10)
	processUnit := GetProcessUnit(sumProcessNumber, len(deadcodes))
	for _, simpleTip := range deadcodes {
		deadCodeTips := strings.Split(simpleTip, ":")
		if len(deadCodeTips) == 4 {
			packageName := PackageNameFromGoPath(deadCodeTips[0])
			line, _ := strconv.Atoi(deadCodeTips[1])
			erroru := Error{
				LineNumber:  line,
				ErrorString: AbsPath(deadCodeTips[0]) + ":" + strings.Join(deadCodeTips[1:], ":"),
			}
			summaries.Lock()
			if summary, ok := summaries.Summaries[packageName]; ok {
				summary.Errors = append(summary.Errors, erroru)
				summaries.Summaries[packageName] = summary
			} else {
				summarie := Summary{
					Name:   PackageAbsPathExceptSuffix(deadCodeTips[0]),
					Errors: make([]Error, 0),
				}
				summarie.Errors = append(summarie.Errors, erroru)
				summaries.Summaries[packageName] = summarie
			}
			summaries.Unlock()
		}
		if sumProcessNumber > 0 {
			s.Sync.LintersProcessChans <- processUnit
			sumProcessNumber = sumProcessNumber - processUnit
		}
	}
	return
}

func (s *StrategyDeadCode) Percentage(summaries Summaries) float64 {
	summaries.Lock()
	defer summaries.Unlock()
	return CountPercentage(len(summaries.Summaries))
}

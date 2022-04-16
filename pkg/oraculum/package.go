package oraculum

type Package struct {
	Data struct {
		BZs []struct {
			BugID            int      `json:"bug_id"`
			Comments         int      `json:"comments"`
			Keywords         []string `json:"keywords"`
			Modified         string   `json:"modified"`
			Priority         string   `json:"priority"`
			PrioritySeverity string   `json:"priority_severity"`
			Private          bool     `json:"private"`
			Release          string   `json:"release"`
			Reported         string   `json:"reported"`
			Severity         string   `json:"severity"`
			Status           string   `json:"status"`
			Title            string   `json:"title"`
			URL              string   `json:"url"`
		} `json:"bzs"`
		Koschei []struct {
			LastSuccess struct {
				Time string `json:"time"`
				URL  string `json:"url"`
			} `json:"last_success"`
			Release string `json:"release"`
			Status  string `json:"status"`
			URL     string `json:"URL"`
		} `json:"koschei"`
	} `json:"data"`
}

func (p *Package) IsFTBFS() bool {
	for _, bz := range p.Data.BZs {
		for _, keyword := range bz.Keywords {
			if keyword == "FTBFS" {
				return true
			}
		}
	}
	return false
}

func (p *Package) IsFTBFS2() bool {
	for _, k := range p.Data.Koschei {
		if k.Status == "failing" {
			return true
		}
	}
	return false
}

func (p *Package) BranchFTBFS(branch string) bool {
	for _, k := range p.Data.Koschei {
		if k.Release == branch {
			if k.Status == "failing" {
				return true
			}
		}
	}
	return false
}

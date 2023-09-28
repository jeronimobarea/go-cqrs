package elasticsearch

type (
	Search struct {
		Query searchQuery `json:"query"`
	}

	searchQuery struct {
		MultiMatch multiMatch `json:"multi_match"`
	}

	multiMatch struct {
		Query            string   `json:"query"`
		Fields           []string `json:"fields"`
		Fuzziness        int      `json:"fuziness"`
		CuttoffFrequency float32  `json:"cutoff_frequency"`
	}
)

func NewSearchWithDefaults(query string, fields []string) Search {
	return newSearch(query, fields, 3, 0.0001)
}

func newSearch(
	query string,
	fields []string,
	fuziness int,
	cutoffFrequency float32,
) Search {
	return Search{
		Query: searchQuery{
			MultiMatch: multiMatch{
				Query:            query,
				Fields:           fields,
				Fuzziness:        fuziness,
				CuttoffFrequency: cutoffFrequency,
			},
		},
	}
}

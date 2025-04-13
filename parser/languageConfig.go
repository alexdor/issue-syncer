package parser

var commentMarkersWithEnds = map[string]CommentMarker{
	".go": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".py": {
		SingleLine: []string{"#"},
		MultiLine: []MultiLineMarker{
			{Start: "'''", End: "'''"},
			{Start: `"""`, End: `"""`},
		},
	},
	".js": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".java": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".c": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".cpp": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".rb": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{{Start: "=begin", End: "=end"}},
	},
	".php": {
		SingleLine: []string{"//", "#"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".rs": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".zig": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{},
	},
	".html": {
		SingleLine: []string{},
		MultiLine:  []MultiLineMarker{{Start: "<!--", End: "-->"}},
	},
	".css": {
		SingleLine: []string{},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".sass": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".less": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".scss": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".ts": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".jsx": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".tsx": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".cs": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".kt": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".pl": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{{Start: "=pod", End: "=cut"}},
	},
	".ex": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{{Start: "@doc \"\"\"", End: "\"\"\""}},
	},
	".exs": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{{Start: "@doc \"\"\"", End: "\"\"\""}},
	},
	".erl": {
		SingleLine: []string{"%"},
		MultiLine:  []MultiLineMarker{},
	},
	".hs": {
		SingleLine: []string{"--"},
		MultiLine:  []MultiLineMarker{{Start: "{-", End: "-}"}},
	},
	".r": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{},
	},
	".vb": {
		SingleLine: []string{"'"},
		MultiLine:  []MultiLineMarker{},
	},
	".swift": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".svelte": {
		SingleLine: []string{},
		MultiLine: []MultiLineMarker{
			{Start: "<!--", End: "-->"},
			{Start: "/*", End: "*/"},
		},
	},
	".scala": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".sh": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{},
	},
	".bash": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{},
	},
	".zsh": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{},
	},
	".yaml": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{},
	},
	".yml": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{},
	},
	".vue": {
		SingleLine: []string{"//"},
		MultiLine: []MultiLineMarker{
			{Start: "<!--", End: "-->"},
			{Start: "/*", End: "*/"},
		},
	},
	"makefile": {
		SingleLine: []string{"#"},
		MultiLine:  []MultiLineMarker{},
	},
	".h": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".hpp": {
		SingleLine: []string{"//"},
		MultiLine:  []MultiLineMarker{{Start: "/*", End: "*/"}},
	},
	".lua": {
		SingleLine: []string{"--"},
		MultiLine: []MultiLineMarker{
			{Start: "--[[", End: "]]"},
		},
	},
}

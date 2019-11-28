package output

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/qneyrat/graphql-coverage/internal/coverage"
)

const templateText = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<style>
			body {
				background: black;
				color: rgb(80, 80, 80);
			}
			body, pre, #legend span {
				font-family: Menlo, monospace;
				font-weight: bold;
			}
			#topbar {
				background: black;
				position: fixed;
				top: 0; left: 0; right: 0;
				height: 42px;
				border-bottom: 1px solid rgb(80, 80, 80);
			}
			#content {
				margin-top: 50px;
			}
			#nav, #legend {
				float: left;
				margin-left: 10px;
			}
			#legend {
				margin-top: 12px;
			}
			#nav {
				margin-top: 10px;
			}
			#legend span {
				margin: 0 5px;
			}
			.cov0 { color: rgb(192, 0, 0) }
			.cov1 { color: rgb(128, 128, 128) }
			.cov2 { color: rgb(116, 140, 131) }
			.cov3 { color: rgb(104, 152, 134) }
			.cov4 { color: rgb(92, 164, 137) }
			.cov5 { color: rgb(80, 176, 140) }
			.cov6 { color: rgb(68, 188, 143) }
			.cov7 { color: rgb(56, 200, 146) }
			.cov8 { color: rgb(44, 212, 149) }
			.cov9 { color: rgb(32, 224, 152) }
			.cov10 { color: rgb(20, 236, 155) }
		</style>
	</head>
	<body>
		<div id="topbar">
			<div id="legend">
				<span>not tracked</span>
				<span class="cov0">no coverage</span>
				<span class="cov1">low coverage</span>
				<span class="cov2">*</span>
				<span class="cov3">*</span>
				<span class="cov4">*</span>
				<span class="cov5">*</span>
				<span class="cov6">*</span>
				<span class="cov7">*</span>
				<span class="cov8">*</span>
				<span class="cov9">*</span>
				<span class="cov10">high coverage</span>
			</div>
		</div>
		<div id="content">
		    <pre>
{{range .CoverFile}} {{span . }} {{end}}
           </pre>
       </div>
   </body>
</html>
`

func Output(filename string, data coverage.WrappedCoverFile) error {
	funcMap := template.FuncMap{
		"span": func(coverLine coverage.CoverLine) string {
			if coverLine.Count <= 0 && (
				strings.HasPrefix(strings.Replace(coverLine.Text, " ", "", -1), "#") ||
					strings.HasPrefix(coverLine.Text, "type") ||
					strings.HasPrefix(coverLine.Text, "}")) {
				return coverLine.Text+"\n"
			}

			if coverLine.Count > 0 {
				return fmt.Sprintf("<span class=\"cov10\" title=\"1\">%s</span>\n", coverLine.Text)
			}

			return fmt.Sprintf("<span class=\"cov0\" title=\"0\">%s</span>\n", coverLine.Text)
		},
	}

	tmpl, err := template.New("coverageHTML").Funcs(funcMap).Parse(templateText)
	if err != nil {
		return err
	}

	output, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}

	return tmpl.Execute(output, data)
}

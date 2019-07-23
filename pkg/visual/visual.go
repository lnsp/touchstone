package visual

import (
	"encoding/json"
	"io"
	"os"
	"text/template"

	"github.com/lnsp/touchstone/pkg/benchmark"
)

var tmpl = template.Must(template.New("").Parse(`
<!doctype html>
<html>

<head>
    <meta charset="utf-8">
    <title>touchstone</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
        integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <script src="https://cdn.jsdelivr.net/npm/chart.js@2.8.0/dist/Chart.min.js"></script>
</head>

<body>
    <header class="container mt-3">
        <div class="row">
            <div class="col">
        		<h1>Touchstone</h1>
            </div>
        </div>
        <hr>
    </header>
    <main class="container">
        <div id="data"></div>
        <script>
			function sort(arr) {
				return arr.sort(function(a, b) { return a - b; });
			}
			function median(arr) {
				arr = sort(arr);
				let size = arr.length;
				let median = Math.floor(size / 2);
				if (size == 1) return arr[0];
				else if (size % 2 == 0) {
					return (arr[median-1] + arr[median]) / 2;
				} else {
					return arr[median];
				}
			}

            let dataContainer = document.getElementById("data");
            let datasets = {{ .Datasets }};
            let indices = {{ .Indices }};

			let colors = {
    			'containerd/runc': 'rgba(30,136,229,0.5)',
    			'containerd/runsc': 'rgba(103,58,183,0.5)',
    			'crio/runc': 'rgba(216,27,96,0.5)',
    			'crio/runsc': 'rgba(244,67,54,0.5)',
			};

            for (op of datasets) {
                console.log("Indexing over " + op.cri + "/" + op.oci);
                for (result of op.results) {
					// Generate index for median computation
					let valueGroups = {};
					for (label of indices[result.name].labels) {
						valueGroups[label] = [];
					}
                    console.log(valueGroups);
					// Aggregate values and sort them
					for (report of result.reports) {
						for (label in report) {
							valueGroups[label].push(report[label]);
						}
					}
					// Compute median
                    let aggregated = [];
                    for (label of indices[result.name].labels) {
                        aggregated.push(median(valueGroups[label]));
                    }
                    indices[result.name].datasets.push({
                        label: op.cri + "/" + op.oci,
                        data: aggregated,
                        borderWidth: 1,
                        backgroundColor: colors[op.cri+'/'+op.oci],
                    });
                }
            }

			Chart.defaults.global.defaultFontFamily = "'American Typewriter', Helvetica, Arial'";
            for (name in indices) {
                var root = document.createElement("p");
                var header = document.createElement("h4");
                header.classList.add("text-monospace");
                var headerDesc = document.createElement("p");
                headerDesc.classList.add("text-muted");
                headerDesc.classList.add("lead");
                headerDesc.appendChild(document.createTextNode(indices[name].description));

                header.appendChild(document.createTextNode(name));
                root.appendChild(header);
                root.appendChild(headerDesc);
                root.appendChild(document.createElement("hr"));
                var canvas = document.createElement("canvas");
                root.appendChild(canvas);
                dataContainer.appendChild(root);

                var ctx = canvas.getContext('2d');
                var myChart = new Chart(ctx, {
                    type: 'bar',
                    data: {
                        labels: indices[name].labels,
                        datasets: indices[name].datasets,
                    },
                    options: {
                        scales: {
                            yAxes: [{
                                ticks: {
                                    beginAtZero: true
                                }
                            }]
                        }
                    }
                });
            }
        </script>
    </main>
</body>

</html>
`))

func HTML(w io.Writer, reportsJSON, indicesJSON []byte) error {
	if err := tmpl.Execute(w, struct {
		Datasets, Indices string
	}{string(reportsJSON), string(indicesJSON)}); err != nil {
		return err
	}
	return nil
}

func Write(name string, entries []benchmark.MatrixEntry, index benchmark.Index) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	entriesBytes, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	indexBytes, err := json.Marshal(index)
	if err != nil {
		return err
	}
	if err := HTML(f, entriesBytes, indexBytes); err != nil {
		return err
	}
	return nil
}

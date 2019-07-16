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
        <h1>touchstone</h1>
        <h3 class="text-muted">benchmarking results</h3>
            </div>
        </div>
        <hr>
    </header>
    <main class="container">
        <div id="data"></div>
        <script>
            let dataContainer = document.getElementById("data");
            let datasets = {{ .Datasets }};
            let indices = {{ .Indices }};

            for (op of datasets) {
                console.log("Indexing over " + op.cri + "/" + op.oci);
                for (result of op.results) {
                    let values = [];
                    for (label of indices[result.name].labels) {
                        values.push(result.aggregated[label]);
                    }
                    indices[result.name].datasets.push({
                        label: op.cri + "/" + op.oci,
                        data: values,
                        borderWidth: 1,
                    });
                    console.log(values);
                }
            }

            for (name in indices) {
                var root = document.createElement("p");
                var header = document.createElement("h4");
                header.appendChild(document.createTextNode(name));
                root.appendChild(header);
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

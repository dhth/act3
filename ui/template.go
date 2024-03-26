package ui

const (
	HTMLTemplText = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
    <title>{{.Title}}</title>
</head>
<body class="bg-slate-900 text-xl">
    <div class="container mx-auto p-4">
        <h1 class="text-stone-50 text-2xl mb-4 font-bold">{{.Title}}</h1>
        <p class="text-stone-300 italic mb-4 mt-8">Generated at {{.Timestamp}}</p>
        <table class="table-auto font-bold text-left">
            <thead>
                <tr class="text-stone-50 bg-slate-700">
                    {{range .Columns -}}
                    <th class="px-4 py-2">{{.}}</th>
                    {{end -}}
                </tr>
            </thead>
            <tbody>
                {{range .Rows -}}
                    <tr>
                    <td class="px-4 py-2 text-blue-400">{{.Key}}</td>
                    {{range .Data -}}
                        {{if .Success}}
                            <td class="px-4 py-2 text-green-400">{{.Result}}</td>
                        {{else}}
                            <td class="px-4 py-2 text-rose-400">{{.Result}}</td>
                        {{end}}
                    {{end -}}
                </tr>
                {{end -}}
            </tbody>
        </table>

        {{if .Failures }}
        <br>
        <hr>
        <br>
        <p class="text-[#fb4934] font-bold italic">Failures</p>
            <br>
            {{range $key, $value := .Failures -}}
            <p class="text-gray-400 italic">{{$key}}: <a class="underline" href="{{$value}}" target="_blank">{{$value}}</a></p>
            <br>
            {{end -}}
        {{end -}}

        {{if .Errors }}
        <br>
        <hr>
        <br>
        <p class="text-red-600 font-bold italic">Errors</p>
            <br>
            {{range $index, $error := .Errors -}}
            <p class="text-gray-400 italic">{{$index}}: {{$error}}</p>
            <br>
            {{end -}}
        {{end -}}
    </div>
</body>
</html>
`
	errorTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
<p>Something went wrong generating the HTML</p>
<p>Error: %s</p>
</body>
</html>
`
)

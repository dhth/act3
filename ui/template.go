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
<body class="bg-[#282828] text-xl">
    <div class="container mx-auto p-4">
        <h1 class="text-[#fbf1c7] text-2xl mb-4 font-bold">{{.Title}}</h1>
        <p class="text-stone-300 italic mb-4 mt-8">Generated at {{.Timestamp}}</p>
        <table class="table-auto w-full text-left">
            <thead>
                <tr class="font-bold text-[#fbf1c7] bg-[#3c3836]">
                    {{range .Columns -}}
                    <th class="px-4 py-2">{{.}}</th>
                    {{end -}}
                </tr>
            </thead>
            <tbody>
                {{range .Rows -}}
                    <tr>
                    <td class="font-bold px-4 py-2 text-[#d3869b]">{{.Key}}</td>
                    {{range .Data -}}
                        {{if .Error}}
                            <td class="px-4 py-2 font-mono">
                                <p class="font-bold inline text-[#fabd2f]">{{.Details.NumberFormatted}}</p>
                                <p class="inline text-[#fabd2f]">{{.Details.Indicator}}</p>
                                <p class="px-1 inline italic text-[#fabd2f]">{{.Details.Context}}</p>
                            </td>
                        {{else if .Success}}
                            <td class="px-4 py-2 font-mono">
                                {{if .Url}}
                                <a class="no-underline" href="{{.Url}}" target="_blank">
                                    <p class="font-bold inline text-[#b8bb26]">{{.Details.NumberFormatted}}</p>
                                    <p class="inline text-[#b8bb26]">{{.Details.Indicator}}</p>
                                    <p class="px-1 inline italic text-[#7c6f64]">{{.Details.Context}}</p>
                                </a>
                                {{else}}
                                <p class="font-bold inline text-[#b8bb26]">{{.Details.NumberFormatted}}</p>
                                <p class="inline text-[#b8bb26]">{{.Details.Indicator}}</p>
                                <p class="px-1 inline italic text-[#7c6f64]">{{.Details.Context}}</p>
                                {{end}}
                            </td>
                        {{else}}
                            <td class="px-4 py-2 font-mono">
                                {{if .Url}}
                                <a class="no-underline" href="{{.Url}}" target="_blank">
                                    <p class="font-bold inline text-[#fb4934]">{{.Details.NumberFormatted}}</p>
                                    <p class="inline text-[#fb4934]">{{.Details.Indicator}}</p>
                                    <p class="px-1 inline italic text-[#7c6f64]">{{.Details.Context}}</p>
                                </a>
                                {{else}}
                                <p class="font-bold inline text-[#fb4934]">{{.Details.NumberFormatted}}</p>
                                <p class="inline text-[#fb4934]">{{.Details.Indicator}}</p>
                                <p class="px-1 inline italic text-[#7c6f64]">{{.Details.Context}}</p>
                                {{end}}
                            </td>
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
        <p class="text-[#fb4934] font-bold italic">Failed Runs</p>
            <br>
        <table class="table-auto text-left">
            {{range $key, $value := .Failures -}}
                <tr>
                    <td class="text-[#a89984] italic px-2 py-1">{{$key}}</td>
                    <td class="text-[#a89984] italic px-2 py-1"><a class="underline" href="{{$value}}" target="_blank">{{$value}}</a></td>
                </tr>
            {{end -}}
        </table>
        {{end -}}

        {{if .Errors }}
        <br>
        <hr>
        <br>
        <p class="text-[#fabd2f] font-bold italic">Errors</p>
            <br>
        <table class="table-auto text-left">
            {{range $index, $error := .Errors -}}
                <tr>
                    <td class="text-[#a89984] italic px-2 py-1">{{$index}}</td>
                    <td class="text-[#a89984] italic px-2 py-1">{{$error}}</td>
                </tr>
            {{end -}}
        </table>
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

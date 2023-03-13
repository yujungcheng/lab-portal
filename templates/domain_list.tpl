{{ define "content" }}

<p>
    <h2>List Domain, total {{ len .Domains }} domain(s).</h2>
</p>

<table>
    <tr>
        <th width=200px>Name</th>
        <th width=320px>UUID</th>
        <th>State</th>
        <th>Vcpu</th>
        <th>Memory</th>
        <th>Disks</th>
        <th>Interfaces</th>
        <th>Actions</th>
    </tr>

    <!-- todo: able to switch multi-line format to show disks and interfaces-->

    {{ range .Domains }}
    <tr>
        <td>{{ .Name }}</td>
        <td>
            <a href="/domain/show/{{.UUID}}">{{ .UUID }}</a>
        </td>
        <td>{{ .StateStr }}</td>
        <td>{{ .Vcpu }}</td>
        <td>{{ .MemoryStr }}</td>

        <td>
            {{ range $key, $value := .Disks }}
                {{ $key }}({{ $value }})
            {{ end }}
        </td>
        <td>
            {{ range $key, $value := .Interfaces }}
                {{ if ne $value.target "" }}
                    <text title="{{ $value.type }} | {{ $value.mac }} | {{ $value.target }}">{{ $value.name }}</text>
                {{ else }}
                    <text title="{{ $value.type }} | {{ $value.mac }}">{{ $value.name }}</text>
                {{ end }}
            {{ end }}
        </td>

        <td>
            {{ if eq .StateStr "Running" }}
                <a href="domain/stop?">ShutOff</a>
            {{ end }}
            {{ if eq .StateStr "Shutoff"}}
                <a href="domain/start">PowerOn</a>
            {{ end }}
        </td>
    </tr>
    {{ end }}

</table>
{{ end }}
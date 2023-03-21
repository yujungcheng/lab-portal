{{ define "content" }}

<p>
    <h3>
        => List All Domains by group : {{ len .DomainsByGroup }} groups(s) |
        <a href="/domains/list-by-group?mode=group">Group By Group Name</a> |
        <a href="/domains/list-by-group?mode=storage">Group By Storage Pool</a> |
        <a href="/domains/list-by-group?mode=network">Group By Network</a> 
    </h3>
</p>

<table>
    <tr>
        <th width=200px>Name</th>
        <th>State</th>
        <th>Vcpu</th>
        <th>Memory</th>
        <th>Disks</th>
        <th>Interfaces</th>
        <th>Actions</th>
    </tr>

    {{ range $group, $domains := .DomainsByGroup }}
    <tr>
        <td colspan=7><h4>{{ $group }} ... {{ len $domains }} domains</h4></td>
    </tr>
    {{ range $domains }}

    <tr>
        <td>
          <text title="{{ .UUID }}">
             <a href="/domains/show?uuid={{ .UUID }}">{{ .Name }}</a>
          </text>

        </td>
        <td>{{ .StateStr }}</td>
        <td>{{ .Vcpu }}</td>
        <td>{{ .MemoryStr }}</td>

        <td>
            {{ range $key, $value := .Disks }}
                <text title="{{ $value.file }}">{{ $value.name }}({{ $value.capacity }})</text>
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
                <a href="domains/stop">ShutOff</a>
            {{ end }}
            {{ if eq .StateStr "Shutoff"}}
                <a href="domains/start">PowerOn</a>
            {{ end }}
        </td>
    </tr>
    {{ end }}



    {{ end }}

</table>

{{ end }}
{{ define "content" }}
<div>

<p>
    <h3>
        => List Cloned Domains ... {{ len . }} domain(s).
    </h3>
</p>

<table class="domain_list">
  <tr>
    <th>Domain Name</th>
    <th>Clone Status</th>
  </tr>

{{ range $key, $value := . }}
  <tr>
    <td>
        <b>{{ $key }}</b>
    </td>
    <td>
         {{ $value }}
    </td>
  </tr>
{{ end }}

</table>

</div>
{{ end }}
{{ define "content" }}

<div>
    <p>
        <td><h3>=> Show Domain ... {{ .Name }}</h3></td>
    </p>

    <form method="POST" action="/domains/update?uuid={{ .UUID }}">
    <table>
        <tr>
            <td><label>UUID</label></td>
            <td>{{ .UUID }}</td>
        </tr>
        <tr>
            <td><label>Name</label></td>
            <td>
                <input type="hidden" name="domainName" value="{{ .Name }}">
                <input type="text" name="newDomainName" value="{{ .Name }}">
            </td>
        </tr>
        <tr>
            <td><label>Status</label></td>
            <td>{{ .StateStr }}</td>
        </tr>
        <tr>
            <td><label>vCPU</label></td>
            <input type="hidden" name="domainVcpu" value="{{ .Vcpu }}">
            <td><input type="text" name="newDomainVcpu" value="{{ .Vcpu }}"></td>
        </tr>
        <tr>
            <td><label>Memory</label></td>
            <input type="hidden" name="domainMem" value="{{ .MemoryStr }}">
            <td><input type="text" name="newDomainMem" value="{{ .MemoryStr }}">GB</td>
        </tr>
        <tr>
            <td><label>Disks</label></td>
            <td>
                <table>
                    <tr>
                        <th>Device Name</th>
                        <th>File Path</th>
                        <th>Allocation</th>
                        <th>Capacity</th>
                        <th>Physical</th>
                    </tr>
                    {{ range .Disks }}
                    <tr>
                        <td>{{ .name }}</td>
                        <td>{{ .file }}</td>
                        <td>{{ .allocation }}GB</td>
                        <td>{{ .capacity }}GB</td>
                        <td>{{ .physical }}GB</td>
                    </tr>
                    {{ end }}
                </table>  
            </td>
        </tr>
        <tr>
            <td><label>Interfaces</label></td>
            <td>
                <table>
                    <tr>
                        <th>Type</th>
                        <th>Name</th>
                        <th>MAC Address</th>
                        <th>Target Name</th>
                    </tr>
                    {{ range .Interfaces }}
                    <tr>
                        <td>{{ .type }} </td>
                        <td>
                        {{ .name }}
                        </td>
                        <td>{{ .mac }}</td>
                        <td>{{ .target }}</td>
                    </tr>
                    {{ end }}
                </table>
            </td>
        </tr>
    </table>

    <br/>
    <center>
        <input name="domainUpdate" type="submit" value="Update Domain"/>
        
    </center>

    </form>

    <br/>

    <!-- ***********************************
        Backup / Delete Domain
    ************************************ -->
    <div style="display: flex; text-align: center;">
    <form method="POST" action="/domains/backup?uuid={{ .UUID }}">
        <label>Enter backup name:</label>
        <input name="domainBackupName" type="text"/>
        <input name="domainBackup" type="submit" value="Backup Domain"/>

    </form>
    &nbsp;|&nbsp;
    <form method="POST" action="/domains/delete?uuid={{ .UUID }}">
        <label>Enter domain name to delete</label>
        <input name="domainDeleteName" type="text"/>
        <input name="domainDelete" type="submit" value="Delete Domain"/>

    </form>
    </div>
</div>

{{ end }}
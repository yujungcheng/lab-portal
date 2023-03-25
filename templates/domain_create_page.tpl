{{ define "content" }}

<div>
    <p><h1>Create New Domain(s)</h1></p>
    <form method="post" action="/domain/create">

        <label>[ Group Domain One ]</label>
        <br/><br/>

        <table>
            <tr>
                <td colspan=4>
                    <label>Name:</label>
                    <input name="g1-name" type="text" value="" />
                    <hr/>
                </td>
            </tr>
            <tr>
                <td>
                    <label><b>OS Type:</b></label>
                    <select id="g1-os-type" name="g1-os-type">
                        <option value="linux">Linux</option>
                        <option value="windows">Windows</option>
                    </select>
                </td>
                <td>
                    <label><b>Number of Domains:</b></abel>
                    <select id="g1-count" name="g1-count" >
                        <option value="1" selected="selected">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="5">&nbsp;5&nbsp;</option>
                        <option value="6">&nbsp;6&nbsp;</option>
                    </select>
                </td>
                <td>
                    <label><b>Number of vCPU:</b></label>
                    <select id="g1-vcpu" name="g1-vcpu" >
                        <option value="1" selected="selected">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                    </select>&nbsp;&nbsp;
                </td>
                <td>
                    <label><b>Memory Size:</b></label>
                    <select id="g1-memory" name="g1-memory" >
                        <option value="1" selected="selected">&nbsp;1GB&nbsp;</option>
                        <option value="2">&nbsp;2GB&nbsp;</option>
                        <option value="4">&nbsp;4GB&nbsp;</option>
                        <option value="8">&nbsp;8GB&nbsp;</option>
                        <option value="16">&nbsp;16GB&nbsp;</option>
                    </select>
                </td>
            </tr>
            <tr>
                <td>
                </td>
            </tr>
        </table>

        <br/>

        <table>
            <tr>
                <td width="40%">
                    <label>Disk Bus:</label>
                    <select id="g1-disk-bus" name="g1-disk-bus" >
                        <option value="virtio" selected="selected">virtio</option>
                        <option value="ide">IDE</option>
                        <option value="sata">SATA</option>
                    </select>
                </td>
                <td>
                    <label>Storage Pool:</label>
                    <select id="g1-pool" name="g1-pool">
                    {{ range $k, $v := .StoragePools }}
                        {{ if eq $v.Name "default" }}
                            <option value="{{ $v.Name }}" selected="selected">{{ $v.Name }}</option>
                        {{ else }}
                            <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                        {{ end }}
                    {{ end }}
                    </select>
                </td>
            </tr>

            <tr>
                <td>
                    <label>Boot Volume Source:</label>
                </td>

                <td>
                    <label>Data Volume Size(GB)</label>

                </td>
            </tr>
        </table>



    </form>

</div>

{{ end }}
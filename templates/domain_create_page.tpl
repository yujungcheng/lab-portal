{{ define "content" }}

<div>
    <p><h1>Create New Domain(s)</h1></p>
    <form method="post" action="/domain/create">

        <label>[ Group Domain One ]</label>
        <br/><br/>

        <table>
            <tr>
                <td>
                    <label>Prefix:</label>
                    <input name="gp1-prefix" type="text" value="" />
                    <label>Name:</label>
                    <input name="gp1-name" type="text" value="" />
                    <label><b>vCPU:</b></label>
                    <select id="gp1-vcpu" name="gp1-vcpu" >
                        <option value="1" selected="selected">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                    </select>&nbsp;&nbsp;
                    <label><b>RAM:</b></label>
                    <select id="gp1-ram" name="gp1-ram" >
                        <option value="1" selected="selected">&nbsp;1GB&nbsp;</option>
                        <option value="2">&nbsp;2GB&nbsp;</option>
                        <option value="4">&nbsp;4GB&nbsp;</option>
                        <option value="8">&nbsp;8GB&nbsp;</option>
                        <option value="16">&nbsp;16GB&nbsp;</option>
                    </select>
                    <label><b>Count:</b></abel>
                    <select id="gp1-count" name="gp1-count" >
                        <option value="1" selected="selected">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="5">&nbsp;5&nbsp;</option>
                        <option value="6">&nbsp;6&nbsp;</option>
                    </select>
                </td>
            </tr>
        </table>

        <br/>

        <table>
            <tr>
                <td>
                    <label>Disk Bus:</label>
                    <select id="gp1-disk-bus" name="gp1-disk-bus" >
                        <option value="virtio" selected="selected">virtio</option>
                        <option value="ide">IDE</option>
                        <option value="sata">SATA</option>
                    </select>
                    <label>Storage Pool:</label>
                    <select id="gp1-pool" name="gp1-pool">
                    {{ range $k, $v := .StoragePools }}
                        {{ if eq $v.Name "default" }}
                            <option value="{{ $v.Name }}" selected="selected">{{ $v.Name }}</option>
                        {{ else }}
                            <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                        {{ end }}
                    {{ end }}
                    </select>
                </td>
                <td>

                </td>
            </tr>

            <tr>
                <td>
                    <label>Boot Disk Source:</label>
                    <select id="gp1-boot-vol" name="gp1-boot-disk">
                    {{ range $k, $v := .Templates }}
                        <option value="{{ $v.UUID }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>Disk#2:</label>
                    <select id="gp1-vol2-size" name="gp1-disk2-size" >
                        <option value="0" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                    &nbsp;|&nbsp;
                    <label>Disk#3:</label>
                    <select id="gp1-vol2-size" name="gp1-disk3-size" >
                        <option value="0" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                    &nbsp;|&nbsp;
                    <label>Disk#4:</label>
                    <select id="gp1-vol2-size" name="gp1-disk4-size" >
                        <option value="0" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB

                </td>

                <td>
                </td>
            </tr>
        </table>

        <br />

        <table>
            <tr>
                <td>
                    <label>Interface Model:</label>
                    <select id="gp1-net-driver" name="gp1-net-driver" >
                        <option value="virtio" selected="selected">virtio</option>
                        <option value="rtl8139">rtl8139</option>
                        <option value="e1000">e1000</option>
                    </select>
                    &nbsp;|&nbsp;
                    <label>NIC#1:</label>
                    <select id="gp1-nic1" name="gp1-nic1">
                    {{ range $k, $v := .Networks }}
                        {{ if eq $v.Name "default" }}
                        <option value="{{ $v.Name }}" selected="selected">{{ $v.Name }}</option>
                        {{ else }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                        {{ end }}
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>NIC#2:</label>
                    <select id="gp1-nic2" name="gp1-nic2">
                    <option value="none" selected="selected"></option>
                    {{ range $k, $v := .Networks }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>Nic#3:</label>
                    <select id="gp1-nic3" name="gp1-nic3">
                    <option value="none" selected="selected"></option>
                    {{ range $k, $v := .Networks }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>

                </td>
                <td>

                </td>
            </tr>
        </table>

        <hr/>
        <label>[ Group Domain Two ]</label>


        <hr/>
        <input type="submit" value="Submit" />

    </form>

</div>

{{ end }}
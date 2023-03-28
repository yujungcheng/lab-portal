{{ define "content" }}

<div>
    <p><h1>Create New Domain(s)</h1></p>
    <form method="POST" action="/domains/create">

        <label>[ Group Domain One ]</label>
        <br/><br/>

        <table>
            <tr>
                <td>
                    <label>Prefix:</label>
                    <input name="group1-prefix" type="text" value="" />
                    <label>Name:</label>
                    <input name="group1-name" type="text" value="" />
                    <label><b>vCPU:</b></label>
                    <select id="group1-vcpu" name="group1-vcpu" >
                        <option value="1" selected="selected">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                    </select>&nbsp;&nbsp;
                    <label><b>RAM:</b></label>
                    <select id="group1-ram" name="group1-ram" >
                        <option value="1" selected="selected">&nbsp;1GB&nbsp;</option>
                        <option value="2">&nbsp;2GB&nbsp;</option>
                        <option value="4">&nbsp;4GB&nbsp;</option>
                        <option value="8">&nbsp;8GB&nbsp;</option>
                        <option value="16">&nbsp;16GB&nbsp;</option>
                    </select>
                    <label><b>Count:</b></abel>
                    <select id="group1-count" name="group1-count" >
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
                    <select id="group1-disk-bus" name="group1-disk-bus" >
                        <option value="virtio" selected="selected">virtio</option>
                        <option value="ide">IDE</option>
                        <option value="sata">SATA</option>
                    </select>
                    <label>Storage Pool:</label>
                    <select id="group1-storage-pool" name="group1-storage-pool">
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
                    <label>Source Boot Disk Domain:</label>
                    <select id="group1-boot-disk-domain" name="group1-boot-disk-domain">
                    {{ range $k, $v := .Templates }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>Disk#2:</label>
                    <select id="group1-disk2-size" name="group1-disk2-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                    &nbsp;|&nbsp;
                    <label>Disk#3:</label>
                    <select id="group1-disk2-size" name="group1-disk3-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                    &nbsp;|&nbsp;
                    <label>Disk#4:</label>
                    <select id="group1-disk2-size" name="group1-disk4-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
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
                    <select id="group1-nic-driver" name="group1-nic-driver" >
                        <option value="virtio" selected="selected">virtio</option>
                        <option value="rtl8139">rtl8139</option>
                        <option value="e1000">e1000</option>
                    </select>
                    &nbsp;|&nbsp;
                    <label>NIC#1:</label>
                    <select id="group1-nic1" name="group1-nic1">
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
                    <select id="group1-nic2" name="group1-nic2">
                    <option value="" selected="selected"></option>
                    {{ range $k, $v := .Networks }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>Nic#3:</label>
                    <select id="group1-nic3" name="group1-nic3">
                    <option value="" selected="selected"></option>
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
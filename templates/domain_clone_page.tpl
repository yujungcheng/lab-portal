{{ define "content" }}

<div>
    <p>
        <h3>=> Clone Domain(s)</h3>
    </p>
    <form method="POST" action="/domains/clone">

        <label>[ Group 1 ]</label>
        <br/><br/>

        <table>
            <tr>
                <td>
                    <label>Source Domain:</label>
                    <select name="group1-original-domain-name">
                    {{ range $k, $v := .Templates }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>New Domain Name Prefix:</label>
                    <input name="group1-prefix" type="text" value="" />
                    <label>Name:</label>
                    <input name="group1-name" type="text" value="" />
                </td>
            </tr>
        </table>
        <br />
        <table>
            <tr>
                <td>
                    <label><b>Number of New Domain:</b></abel>
                    <select name="group1-count" >
                        <option value="1" selected="selected">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="5">&nbsp;5&nbsp;</option>
                        <option value="6">&nbsp;6&nbsp;</option>
                    </select>
                    &nbsp;|&nbsp;
                    <label><b>Domain vCPU:</b></label>
                    <select name="group1-vcpu" >
                        <option value="" selected="selected"></option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                    </select>
                    &nbsp;&nbsp;
                    <label><b>Domain RAM Size:</b></label>
                    <select name="group1-ram" >
                        <option value="" selected="selected"></option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                    </select>
                    GB
                </td>
                
            </tr>
        </table>

        <br />

        <table>
            <tr>
                <td>
                    <label>Interface Model:</label>
                    <select name="group1-nic-driver" >
                        <option value="virtio" selected="selected">virtio</option>
                        <option value="rtl8139">rtl8139</option>
                        <option value="e1000">e1000</option>
                    </select>
                    &nbsp;|&nbsp;
                    <label>Interface Network #1:</label>
                    <select name="group1-nic1">
                    {{ range $k, $v := .Networks }}
                        {{ if eq $v.Name "default" }}
                        <option value="{{ $v.Name }}" selected="selected">{{ $v.Name }}</option>
                        {{ else }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                        {{ end }}
                    {{ end }}
                    </select>
                    &nbsp;&nbsp;
                    <label>#2:</label>
                    <select name="group1-nic2">
                    <option value="" selected="selected"></option>
                    {{ range $k, $v := .Networks }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>
                    &nbsp;&nbsp;
                    <label>#3:</label>
                    <select name="group1-nic3">
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

        <br/>

        <table>
            <tr>
                <td>
                    <label>Exist Storage Pool:</label>
                    <select name="group1-storage-pool">
                    {{ range $k, $v := .StoragePools }}
                        {{ if eq $v.Name "default" }}
                            <option value="{{ $v.Name }}" selected="selected">{{ $v.Name }}</option>
                        {{ else }}
                            <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                        {{ end }}
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>New Storage Pool Name:</label>
                    <input name="group1-new-storage-pool" type="text" value="" />
                    <label>Path:</label>
                    <input name="group1-new-storage-pool-path" type="text" value="" />
                </td>
            </tr>
            <tr>
                <td>
                    <label>Disk Bus:</label>
                    <select name="group1-disk-bus" >
                        <option value="virtio" selected="selected">VirtIO</option>
                        <option value="ide">IDE</option>
                        <option value="sata">SATA</option>
                    </select>
                    &nbsp;&nbsp;
                    <label>Data Disk #1:</label>
                    <select name="group1-disk1-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                    &nbsp;&nbsp;
                    <label>#2:</label>
                    <select name="group1-disk2-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                    &nbsp;&nbsp;
                    <label>#3:</label>
                    <select name="group1-disk3-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                </td>
            </tr>
        </table>

        <br/>
        <label>[ Group 2 ]</label>
        <br/><br/>

        <table>
            <tr>
                <td>
                    <label>Source Domain:</label>
                    <select name="group2-original-domain-name">
                    {{ range $k, $v := .Templates }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>New Domain Name Prefix:</label>
                    <input name="group2-prefix" type="text" value="" />
                    <label>Name:</label>
                    <input name="group2-name" type="text" value="" />
                </td>
            </tr>
        </table>
        <br />
        <table>
            <tr>
                <td>
                    <label><b>Number of New Domain:</b></abel>
                    <select name="group2-count" >
                        <option value="1" selected="selected">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="5">&nbsp;5&nbsp;</option>
                        <option value="6">&nbsp;6&nbsp;</option>
                    </select>
                    &nbsp;|&nbsp;
                    <label><b>Domain vCPU:</b></label>
                    <select name="group2-vcpu" >
                        <option value="" selected="selected"></option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="3">&nbsp;3&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                    </select>
                    &nbsp;&nbsp;
                    <label><b>Domain RAM Size:</b></label>
                    <select name="group2-ram" >
                        <option value="" selected="selected"></option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                    </select>
                    GB
                </td>
                
            </tr>
        </table>

        <br />

        <table>
            <tr>
                <td>
                    <label>Interface Model:</label>
                    <select name="group2-nic-driver" >
                        <option value="virtio" selected="selected">virtio</option>
                        <option value="rtl8139">rtl8139</option>
                        <option value="e1000">e1000</option>
                    </select>
                    &nbsp;|&nbsp;
                    <label>Interface Network #1:</label>
                    <select name="group2-nic1">
                    {{ range $k, $v := .Networks }}
                        {{ if eq $v.Name "default" }}
                        <option value="{{ $v.Name }}" selected="selected">{{ $v.Name }}</option>
                        {{ else }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                        {{ end }}
                    {{ end }}
                    </select>
                    &nbsp;&nbsp;
                    <label>#2:</label>
                    <select name="group2-nic2">
                    <option value="" selected="selected"></option>
                    {{ range $k, $v := .Networks }}
                        <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                    {{ end }}
                    </select>
                    &nbsp;&nbsp;
                    <label>#3:</label>
                    <select name="group2-nic3">
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

        <br/>

        <table>
            <tr>
                <td>
                    <label>Exist Storage Pool:</label>
                    <select name="group2-storage-pool">
                    {{ range $k, $v := .StoragePools }}
                        {{ if eq $v.Name "default" }}
                            <option value="{{ $v.Name }}" selected="selected">{{ $v.Name }}</option>
                        {{ else }}
                            <option value="{{ $v.Name }}">{{ $v.Name }}</option>
                        {{ end }}
                    {{ end }}
                    </select>
                    &nbsp;|&nbsp;
                    <label>New Storage Pool Name:</label>
                    <input name="group2-new-storage-pool" type="text" value="" />
                    <label>Path:</label>
                    <input name="group2-new-storage-pool-path" type="text" value="" />
                </td>
            </tr>
            <tr>
                <td>
                    <label>Disk Bus:</label>
                    <select name="group2-disk-bus" >
                        <option value="virtio" selected="selected">VirtIO</option>
                        <option value="ide">IDE</option>
                        <option value="sata">SATA</option>
                    </select>
                    &nbsp;&nbsp;
                    <label>Data Disk #1:</label>
                    <select name="group2-disk1-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                    &nbsp;&nbsp;
                    <label>#2:</label>
                    <select name="group2-disk2-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                    &nbsp;&nbsp;
                    <label>#3:</label>
                    <select name="group2-disk3-size" >
                        <option value="" selected="selected">&nbsp;&nbsp;</option>
                        <option value="1">&nbsp;1&nbsp;</option>
                        <option value="2">&nbsp;2&nbsp;</option>
                        <option value="4">&nbsp;4&nbsp;</option>
                        <option value="8">&nbsp;8&nbsp;</option>
                        <option value="16">&nbsp;16&nbsp;</option>
                        <option value="32">&nbsp;32&nbsp;</option>
                    </select> GB
                </td>
            </tr>
        </table>

        <br/>
        <center>
            <input type="submit" value="Submit" />
        </center>
    </form>

</div>

{{ end }}
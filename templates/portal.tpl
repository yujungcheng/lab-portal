
{{ define "title" }}Virtual Lab Portal{{ end }}

{{ define "style" }}

    label {
        font-weight: bold;
    }

    .vertical {
        border-left: 2px solid;
        height: 100px;
        position: absolute;
        left:30%;
    }

    /* table style */
    
    table, tr, td {
        border:none !important;
    }

    table {
        font-family: arial, sans-serif;
        border-collapse: collapse;
        width: 100%;
    }
    td, th {
        border: 1px solid #dddddd;
        text-align: left;
        padding: 4px;
        vertical-align: top;
    }
    /*
    tr:nth-child(even) {
        background-color: #dddddd;
    }
    */
    /*
    tbody tr:hover {
        background-color: #dddddd;
    }
    */
    th {
        background-color: #dddddd;
    }

    .domain_list tr:hover {
        background-color: #dddddd;
    }

    /* dropdown menu */
    div.menu {
        font-weight: bolder;
        normal, bold, bolder, lighter, 100~900 
        font-weight: 900;
    }
    div.action {
        padding-left: 140px;
    }

    .dropbtn {
        padding: 8px;
        font-size: 16px;
        border: none;
    }

    .dropdown {
        position: relative;
        display: inline-block;
    }

    .dropdown-content {
        display: none;
        position: absolute;
        background-color: LightGray;
        min-width: 200px;
        box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
        z-index: 1;
    }

    .dropdown-content a {
        color: black;
        padding: 6px 8px;
        text-decoration: none;
        display: block;
    }

    .dropdown-content a:hover {background-color: Lavender;}
    .dropdown:hover .dropdown-content {display: block;}
    .dropdown:hover .dropbtn {background-color: Lavender;}

    /* items https://www.w3schools.com/colors/colors_names.asp */
    .active-item {
        //border-style:solid;
        //border-width: thin;
        background-color: LawnGreen;
    }
    .inactive-item {
        background-color: Coral;
    }
    .list-item {
        background-color: Orange;
    }
{{ end }}


{{ define "menu" }}
    {{/* comment */}}
    {{- /* comment with white space trimmed */ -}}
    {{- /* comment with white space and newline trimmed */ -}}
    <div class="menu">

        <div class="dropdown">
            <button class="dropbtn">
                <a href="/domains">Domains</a>
            </button> |
            <div class="dropdown-content">
                <a href="/domains/create-page">[ Create Domain ]</a>
                <a href="/domains/delete-page">[ Delete Domain ]</a>
                <a href="/domains/update-page">[ Update Domain ]</a>
                <a href="/domains/backup-page">[ Backup Domain ]</a>
            </div>
        </div>


    </div>
{{ end }}

{{ define "footer" }}
    <b>This is footer...</b>
{{ end }}



{{ template "base" . }}

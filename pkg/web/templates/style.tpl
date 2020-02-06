body {
    background-color: #FFFFFF;
    background-image: url("{{url}}/static/img/back-chateau.jpg");
    background-repeat: no-repeat;
    background-position: right;
    background-attachment: fixed;
    background-size: contain;
}

body.troops {
    background-image: url("{{url}}/static/img/back-general.jpg");
}

body.buildings {
    background-image: url("{{url}}/static/img/back-magie.jpg");
}

body.science {
    background-image: url("{{url}}/static/img/back-science.jpg");
}

* {
    font-family: Helvetica, Tahoma, sans-serif;
    color: black;
}

main h1 {
    font-size: 36px;
}

header h1 {
    font-size: 36px;
    color: deeppink;
}

header h2 {
    font-size: 24px;
    color: darkorange;
}

h2 {
    font-size: 20px;
}

h3 {
    font-size: 16px;
}

p, li, td, div, td {
    font-size: 14px;
}

a {
    color: seagreen;
}

em {
    color: #000000;
    font-weight: bold;
}

td {
    padding: 0px 5px 0px 5px;
    text-align: right;
}

thead td {
    border-bottom: 1px solid black;
}

thead td {
    font-weight: bold;
}

tfoot td {
    font-weight: bold;
    border-top: 1px solid black;
}

td.title {
    padding: 0px 20px 0px 5px;
}

footer p {
    text-align: center;
}

footer a {
    color: yellowgreen;
}

aside p {
    color: orange;
    font-weight: bold;
    font-size: 24px;
}

nav a {
    color: yellowgreen;
    font-weight: bolder;
    font-size: 20px;
    text-decoration: none;
}

div {
    background-color: rgba(243, 243, 240, 0.6);
}

div h2 {
    /*line-height: 20px;*/
    padding: 5px;
    margin: 0px;
    color: darkslategray;
}


@media screen {
    @media (orientation: landscape) and (min-width: 1300px) {
        body {
            display: grid;
            grid-template-columns: 200px 1024px auto;
            align-content: flex-start;
        }

        header {
            grid-column: 1 / 4;
            grid-row: 1;
            display: block;
            margin-top: 10px;
            margin-bottom: 30px;
        }

        header h1, header h2 {
            padding-top: 0px;
            padding-bottom: 0px;
            margin-top: 0px;
            margin-bottom: 0px;
        }

        nav {
            grid-column: 1;
            grid-row: 2;
            display: block;
            vertical-align: top;
        }

        aside {
            grid-column: 3;
            grid-row: 2;
            display: block;
        }

        footer {
            grid-column: 1 / 4;
            grid-row: 4;
            width: 100%;
        }

        nav a {
            display: block;
        }

        main {
            grid-column: 2;
            grid-row: 2 / 4;

            display: grid;
            grid-template-columns: 50% 50%;
        }

        main div {
            border: 1px solid darkslategray;
            margin: 5px;
            padding: 10px;
        }

        main div.large {
            grid-column: 1 / 3;
        }
    }
    @media (orientation: landscape) and (max-width: 1300px) and (min-width: 800px) {
        body {
            display: grid;
            grid-template-columns: 200px 600px;
            align-content: flex-start;
        }

        header {
            grid-column: 1 / 4;
            grid-row: 1;
            display: block;
            margin-top: 10px;
            margin-bottom: 30px;
        }

        header h1, header h2 {
            padding-top: 0px;
            padding-bottom: 0px;
            margin-top: 0px;
            margin-bottom: 0px;
        }

        nav {
            grid-column: 1;
            grid-row: 2;
            display: block;
            vertical-align: top;
        }

        aside {
            grid-column: 3;
            grid-row: 2;
            display: block;
        }

        footer {
            grid-column: 1 / 4;
            grid-row: 4;
            width: 100%;
        }

        nav a {
            display: block;
        }

        main {
            grid-column: 2;
            grid-row: 2 / 4;
        }

        main div {
            border: 1px solid darkslategray;
            margin: 5px;
            padding: 10px;
        }

    }
    @media (orientation: portrait) or (max-width: 800px) {
        header h1, header h2 {
            padding-top: 0px;
            padding-bottom: 0px;
            margin-top: 0px;
            margin-bottom: 0px;
        }

        nav a {
            display: inline;
            padding-left: 5px;
            padding-right: 5px;
        }

        main section {
            width: 100%;
            margin: auto;
        }

        section.col2 {
            display: block;
            grid-template-columns: auto;
        }

        section.col2 div {
            border: 2px solid darkslategray;
            margin: 5px;
            padding: 10px;
            width: 460px;
        }

        section.col2 div.large {
            grid-column: 1 / 3;
            width: 960px;
        }

        main div {
            border: 1px solid darkslategray;
            margin: 5px;
            padding: 10px;
        }

    }
}

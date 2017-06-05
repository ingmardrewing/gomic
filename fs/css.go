package fs

type cssGen struct{}

func newCss() *cssGen {
	return &cssGen{}
}

func (cg *cssGen) getCss() string {
	return cssTemplate
}

const cssTemplate = `
.copyright,
header {
	width: 800px;
	margin: 0 auto;
}

ul.archive {
	list-style-type: none;
}

ul.archive li {
	display: inline-block;
	margin: 10px;
}

.copyright{
	margin-top: 30px;
	margin-bottom: 60px;
}

h3 {
	font-family: Arial Black;
	text-align: left;
	text-transform: uppercase;
}

header .home {
    display: block;
    line-height: 80px;
    background: url(https://devabo.de/imgs/header_devabo_de.png) no-repeat 0px -0px;
    height: 30px;
    width: 800px;
    text-align: left;
    color: #000;
    margin-bottom: 0px;
	margin-top: 0;
    background-color: transparent;
}

header .orange {
	display: block;
    height: 2.2em;
    background-color: #FF8800;
    color: #FFFFFF;
    line-height: 1em;
    padding: 0.5em;
    box-sizing: border-box;
	width: 100%;
    font-size: 24px;
	font-family: Arial Black;
	text-transform: uppercase;
    text-decoration: underline;
	margin-bottom: 1rem;
}

body {
	text-align: center;
	margin: 0;
	padding: 0;
	border: 0;
	font-family: Arial, Helvetica, sans-serif;
}

#disqus_thread,
main {
	width: 800px;
	margin: 0 auto;
	text-align: left;
}

main nav {
	text-align: right;
}

footer {
	position: fixed;
	bottom: 0;
	width: 100%;
	text-align: center;
	z-index: 100;
}

footer nav {
	border-top: 1px solid black;
	position: relative;
	background-color: white;
	min-height: 45px;
	width: 800px;
	margin: 0 auto;
}

nav a {
	font-family: Arial Black;
	color: black;
	text-decoration: none;
	height: 100%;
	display: inline-block;
	padding: 10px;
	text-transform: uppercase;
}

.spacer {
	height: 80px;
}

#cookie-law-info-bar {
	font-size: 10pt;
	margin: 0 auto;
	padding: 5px 0;
	position: fixed;
	top: 0;
	left: 0;
	text-align: center;
	width: 100%;
	z-index: 9999;
	background-color: white;
	border: 1px solid black;
}

#cookie-law-info-again {
	font-size: 10pt;
	margin: 0;
	padding: 2px 10px;
	text-align: center;
	z-index: 9999;
	cursor: pointer;
}

#cookie-law-info-bar span {
	vertical-align: middle;
}

/** Buttons (http://papermashup.com/demos/css-buttons) */
.cli-plugin-button, .cli-plugin-button:visited {
	display: inline-block;
	padding: 5px 10px 6px;
	color: #fff;
	background-color: #000;
	text-decoration: none;
	-moz-border-radius: 6px;
	-webkit-border-radius: 6px;
	-moz-box-shadow: 0 1px 3px rgba(0,0,0,0.6);
	-webkit-box-shadow: 0 1px 3px rgba(0,0,0,0.6);
	text-shadow: 0 -1px 1px rgba(0,0,0,0.25);
	border-bottom: 1px solid rgba(0,0,0,0.25);
	position: relative;
	cursor: pointer;
	margin: auto 10px;
}

.cli-plugin-button:hover {
	background-color: #111;
	color: #fff;
}

.cli-plugin-button:active {
	top: 1px;
}

.small.cli-plugin-button, .small.cli-plugin-button:visited {
	font-size: 11px;
}

.cli-plugin-button, .cli-plugin-button:visited,
	.medium.cli-plugin-button, .medium.cli-plugin-button:visited {
	font-size: 13px;
	font-weight: bold;
	line-height: 1;
	text-shadow: 0 -1px 1px rgba(0,0,0,0.25);
}

.large.cli-plugin-button, .large.cli-plugin-button:visited {
	font-size: 14px;
	padding: 8px 14px 9px;
}

.super.cli-plugin-button, .super.cli-plugin-button:visited {
	font-size: 34px;
	padding: 8px 14px 9px;
}

.pink.cli-plugin-button, .magenta.cli-plugin-button:visited {
	background-color: #e22092;
}

.pink.cli-plugin-button:hover {
	background-color: #c81e82;
}

.green.cli-plugin-button, .green.cli-plugin-button:visited {
	background-color: #91bd09;
}

.green.cli-plugin-button:hover {
	background-color: #749a02;
}

.red.cli-plugin-button, .red.cli-plugin-button:visited {
	background-color: #e62727;
}

.red.cli-plugin-button:hover {
	background-color: #cf2525;
}

.orange.cli-plugin-button, .orange.cli-plugin-button:visited {
	background-color: #ff5c00;
}

.orange.cli-plugin-button:hover {
	background-color: #d45500;
}

.blue.cli-plugin-button, .blue.cli-plugin-button:visited {
	background-color: #2981e4;
}

.blue.cli-plugin-button:hover {
	background-color: #2575cf;
}

.yellow.cli-plugin-button, .yellow.cli-plugin-button:visited {
	background-color: #ffb515;
}

.yellow.cli-plugin-button:hover {
	background-color: #fc9200;
}

.nl_container {
	position: fixed;
	left: 0;
	top: 0;
	width: 100%;
	height: 100%;
	background-color: red;
}

.nl_container_hidden {
	display: none;
}
`

package dinghy

import "strings"

var icons []string

func init() {
	icons = []string{"3g2", "3ga", "3gp", "7z", "aac", "aa", "accdb", "accdt", "ace", "ac", "adn", "aifc", "aiff", "aif", "ai", "ait", "amr", "ani", "apk", "applescript", "app", "asax", "asc", "ascx", "asf", "ash", "ashx", "asm", "asmx", "asp", "aspx", "asx", "aup", "au", "avi", "axd", "aze", "bak", "bash", "bat", "bin", "blank", "bmp", "bowerrc", "bpg", "browser", "bz2", "bzempty", "cab", "cad", "caf", "cal", "catalog.json", "cdda", "cd", "cer", "cfg", "cfml", "cfm", "cgi", "chm", "class", "cmd", "codekit", "code-workspace", "coffeelintignore", "coffee", "compile", "com", "config", "conf", "cpp", "cptx", "cr2", "crdownload", "crt", "crypt", "csh", "cson", "csproj", "css", "cs", "c", "csv", "cue", "cur", "dart", "data", "dat", "dbf", "db", "deb", "default", "dgn", "dist", "diz", "dll", "dmg", "dng", "docb", "docm", "doc", "docx", "dotm", "dot", "dotx", "download", "dpj", "dsn", "ds_store", "dtd", "dwg", "dxf", "editorconfig", "elf", "el", "eml", "enc", "eot", "eps", "epub", "eslintignore", "exe", "f4v", "fax", "fb2", "flac", "fla", "flv", "fnt", "folder", "fon", "gadget", "gdp", "gem", "gif", "gitattributes", "gitignore", "go", "gpg", "gpl", "gradle", "gz", "handlebars", "hbs", "heic", "hlp", "hsl", "hs", "h", "html", "htm", "ibooks", "icns", "ico", "ics", "idx", "iff", "ifo", "image", "img", "iml", "inc", "indd", "info", "inf", "ini", "in", "inv", "iso", "j2", "jar", "java", "jpeg", "jpe", "jpg", "json", "jsp", "js", "jsx", "key", "kf8", "kmk", "ksh", "kts", "kt", "kup", "less", "lex", "licx", "lisp", "lit", "lnk", "lock", "log", "lua", "m2v", "m3u8", "m3u", "m4a", "m4r", "m4", "m4v", "map", "master", "mc", "mdb", "mdf", "md", "me", "midi", "mid", "mi", "mk", "mkv", "mm", "mng", "mobi", "mod", "mo", "mov", "mp2", "mp3", "mp4", "mpa", "mpd", "mpeg", "mpe", "mpga", "mpg", "mpp", "mpt", "msg", "msi", "msu", "m", "nef", "nes", "nfo", "nix", "npmignore", "ocx", "odb", "ods", "odt", "ogg", "ogv", "ost", "otf", "ott", "ova", "ovf", "p12", "p7b", "pages", "part", "pcd", "pdb", "pdf", "pem", "pfx", "pgp", "phar", "php", "ph", "pid", "pkg", "plist", "pl", "pm", "png", "pom", "po", "pot", "potx", "pps", "ppsx", "pptm", "ppt", "pptx", "prop", "ps1", "psd", "psp", "ps", "pst", "pub", "pyc", "py", "qt", "ram", "rar", "ra", "raw", "rb", "rdf", "rdl", "reg", "resx", "retry", "rm", "rom", "rpm", "rpt", "rsa", "rss", "rst", "rtf", "rub", "ru", "sass", "scss", "sdf", "sed", "sh", "sitemap", "sit", "skin", "sldm", "sldx", "sln", "sol", "sphinx", "sqlite", "sql", "step", "stl", "svg", "swd", "swf", "swift", "swp", "sys", "tar", "tax", "tcsh", "tex", "tfignore", "tga", "tgz", "tiff", "tif", "tmp", "tmx", "torrent", "tpl", "ts", "tsv", "ttf", "twig", "txt", "udf", "vbproj", "vbs", "vb", "vcd", "vcf", "vcs", "vdi", "vdx", "vmdk", "vob", "vox", "vscodeignore", "vsd", "vss", "vst", "vsx", "vtx", "war", "wav", "wbk", "webinfo", "webm", "webp", "wma", "wmf", "wmv", "woff2", "woff", "wps", "wsf", "xaml", "xcf", "xfl", "xlm", "xlsm", "xls", "xlsx", "xltm", "xlt", "xltx", "xml", "xpi", "xps", "xrb", "xsd", "xsl", "xspf", "xz", "yaml", "yml", "zip", "zsh", "z"}
}

func addIcons(ff []File) {
	for i, f := range ff {
		ff[i].Icon = "blank"

		for _, ext := range icons {
			if strings.HasSuffix(f.Name, "."+ext) {
				ff[i].Icon = ext
				break
			}
		}
	}
}

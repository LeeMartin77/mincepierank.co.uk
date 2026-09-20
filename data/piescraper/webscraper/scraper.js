
let cache = {}

function download_image(url) {
    var urlprsd = URL.parse(url)
    if (!urlprsd) {
        return
    }
    if (cache[urlprsd.origin + urlprsd.pathname]) {
        // we have already downloaded this
        return
    }
    cache[urlprsd.origin + urlprsd.pathname] = true;

    var page = window.location.href
    fetch(url)
        .then(response => response.blob())  
        .then(blob => {
            var a = document.createElement('a');
            a.href = URL.createObjectURL(blob);  
            a.download = page + url; 
            a.style.display = 'none';  
            document.body.appendChild(a);
            a.click(); 
            document.body.removeChild(a); 
        })
        .catch(err => {
            console.error('Error fetching the image:', err);
        });
}

function find_and_download_all_img() {
    const all_images = document.getElementsByTagName("img")
    for (const img of all_images) {
        if (img.src && img.src.startsWith("http") && img.width > 200){
            download_image(img.src)
        }
    }
}

function extracttextrecursively(node) {
    if (node.localName == "script" || node.localName == "link" || node.localName == "style") {
        return ""
    }   
    let text = "";

    for (const child of node.children) {
        text += extracttextrecursively(child);
    }

    if (node.children.length == 0 && node.innerText) {
        text += " " + node.innerText;
    }

    return text;
}

function download_pagetext() {
    var page = window.location.href
    var a = document.createElement('a');
    a.href = URL.createObjectURL(new Blob([extracttextrecursively(document.body)], { type: 'text/plain' }));  
    a.download = page; 
    a.style.display = 'none';  
    document.body.appendChild(a);
    a.click(); 
    document.body.removeChild(a); 
}


function scrape() {
    download_pagetext();
    find_and_download_all_img();
}
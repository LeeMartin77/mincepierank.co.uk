
const fs = require('fs');
const readline = require('readline');

const [,, directory] = process.argv;
if (!directory) {
    console.error('Usage: node organiser.js <directory>');
    process.exit(1);
}

// get all the text files
const files = fs.readdirSync(directory)

const locations = files.filter(f => f.endsWith(".txt")).map(x => x.replace(".txt", ""))
// turn them into urls and use the hostname as a folder
const hostnames = new Set()
locations.forEach((x) => {
    x = x.replace("___", "://") // get it back to http/https
    x = x.replace("_", "/") // bit of an assumption we won't have a domain underscore
    url = URL.parse(x)
    hostnames.add(url.hostname)
})
// find all things with the same prefix as the hostname and shove them into the hostname folder
hostnames.forEach((x) => {
    fs.mkdirSync(directory + "/" + x)
})

files.forEach((f) => {
    hostnames.forEach(hn => {
        if (f.includes(hn)) {
            fs.renameSync(directory + "/" + f, directory + "/"+hn+"/"+f)
        }
    })
})
export function api(path, method="GET", body={}) {
    const base_path = "/api/"

    return new Promise((resolve) => {
        const init = {
            method: method,
            headers: { "Content-Type": "application/json" }
        }
        if(method !== "GET" && body !== {}) {
            init.body = JSON.stringify(body)
        }
        fetch(base_path+path, init).then((resp) => {
            resp.json().then((data) => {
                if(data.status !== "ok") {
                    error("request error "+data.status+": "+data.data)
                    resolve(null)
                }
                resolve(data.data)
            }).catch((err) => {
                console.error(err)
                resolve(null)
            })
        }).catch((err) => {
            error("request error: "+err)
            resolve(null)
        })
    })
}


export function parseHtml(data) {
    let html = (new DOMParser().parseFromString(data, "text/html")).documentElement.textContent;

    let start = 0;
    while(true) {
        let i = html.indexOf("class=\"spoiler\"", start);
        if(i >= 0) {
            i += 15
            let id = "sp_" + Math.floor(Math.random() * (1000 + 1));
            let add = " id=\""+id+"\" onclick=\"document.querySelector('#"+id+"').classList.remove('spoiler')\" ";
            html = html.substring(0,i) + add + html.substring(i)
            start += i + add.length
        } else {
            break;
        }
    }
    return html
}

function error(err) {
    //alert(err)
    console.error(err)
}
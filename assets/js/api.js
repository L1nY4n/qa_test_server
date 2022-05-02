
// 异步 get
async function Get(url) {
    let response = await fetch(url)
    if (response.ok) {
        return await response.json();
    } else {
        alert("HTTP-Error: " + response.status);
    }
}

// post
async function Post(path, data) {
    let response = await fetch(path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8'
        },
        body: JSON.stringify(data)
    })
    if (response.ok) {
        return await response.json();
    } else {
        alert("HTTP-Error: " + response.status);
    }
}

// 获取设备信息
async function DeviceInFo(id, e) {

    e.classList.add("loading")
    const json = await Get(`/device/${id}`)
    e.classList.remove("loading")
    const c = json.success ? "success" : "faild"
    e.classList.add(c)
    setTimeout(() => { e.classList.remove(c) }, 1500)
    return json
}
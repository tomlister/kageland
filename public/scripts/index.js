const res = await fetch("/api/top");
const top = await res.json();

const res2 = await fetch("/api/new");
const newShaders = await res2.json();


const row = document.getElementById("top-shaders");
const row2 = document.getElementById("new-shaders");

// hacky
const waitLoad = (v, key) => {
    const el = document.getElementById(`${key}-viewer-${v.id}`);
    if (el?.contentWindow && el?.contentWindow?.compileShader) {
        document.getElementById(`viewer-${v.id}`).contentWindow.compileShader(v.frag_shader, v.image_1, v.image_2, v.image_3, v.image_4);
        return true;
    }
    return false;
}

newShaders.forEach(v => {
    const el = document.createElement("div");
    el.className = "col-sm-4";
    el.innerHTML = `<div class="card">
        <iframe allowtransparency="true" style="background: #000000;" class="card-img-top" id="new-viewer-${v.id}" src="_viewer.html" width="18rem"></iframe>
        <div class="card-body">
        <h5 class="card-title">${v.name}</h5>
        <div class="react-group" aria-label="Like shader" id="likeButton">
            <span id="likeButtonIcon" class="material-icons">favorite_border</span>
            <div id="likeCount">${v.likes}</div>
        </div>
        <div class="view-group">
            <span class="material-icons">visibility</span>
            <div id="viewCount">${v.views}</div>
        </div>
        <br />
        <a href="/shader?id=${v.id}" class="btn btn-primary">View</a>
        </div>
    </div>`;
    row2.append(el);
    let i = setInterval(() => {
        if (waitLoad(v, "new")) clearInterval(i);
    }, 100)
});

top.forEach(v => {
    const el = document.createElement("div");
    el.className = "col-sm-4";
    el.innerHTML = `<div class="card">
        <iframe allowtransparency="true" style="background: #000000;" class="card-img-top" id="top-viewer-${v.id}" src="_viewer.html" width="18rem"></iframe>
        <div class="card-body">
        <h5 class="card-title">${v.name}</h5>
        <div class="react-group" aria-label="Like shader" id="likeButton">
            <span id="likeButtonIcon" class="material-icons">favorite_border</span>
            <div id="likeCount">${v.likes}</div>
        </div>
        <div class="view-group">
            <span class="material-icons">visibility</span>
            <div id="viewCount">${v.views}</div>
        </div>
        <br />
        <a href="/shader?id=${v.id}" class="btn btn-primary">View</a>
        </div>
    </div>`;
    row.append(el);
    let i = setInterval(() => {
        if (waitLoad(v, "top")) clearInterval(i);
    }, 100)
});
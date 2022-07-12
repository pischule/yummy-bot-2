const img = new Image();

const container = document.getElementById("container");
const canvas = document.getElementById("canvas");
const file = document.getElementById("file");
const ctx = canvas.getContext("2d");

let x = 0;
let y = 0;
let rectangles = [];
let isDrawing = false;
let movingIndex = -1;
let lastPos = null;
let imageLoaded = false;

function rect2canvas(r) {
  return {
    x: r.x * canvas.width,
    y: r.y * canvas.height,
    w: r.w * canvas.width,
    h: r.h * canvas.height,
  };
}

function initImage() {
  const cs = getComputedStyle(container);
  const paddingX = parseFloat(cs.paddingLeft) + parseFloat(cs.paddingRight);
  const borderX =
    parseFloat(cs.borderLeftWidth) + parseFloat(cs.borderRightWidth);
  canvas.width = container.clientWidth - paddingX - borderX;
  const ratio = canvas.width / img.width;
  canvas.height = img.height * ratio;
  imageLoaded = true;
  draw();
}

function chunks(bigarray, size) {
  const arrayOfArrays = [];
  for (let i = 0; i < bigarray.length; i += size) {
    arrayOfArrays.push(bigarray.slice(i, i + size));
  }
  return arrayOfArrays;
}

function draw() {
  ctx.globalCompositeOperation = "copy";
  // prettier-ignore
  ctx.drawImage(img,
    0, 0, img.width, img.height,
    0, 0, canvas.width, canvas.height
  );

  ctx.globalCompositeOperation = "difference";
  ctx.fillStyle = "white";
  for (let r of rectangles) {
    const cr = rect2canvas(r);
    ctx.fillRect(cr.x, cr.y, cr.w, cr.h);
  }

  ctx.globalCompositeOperation = "source-over";
  ctx.font = "50px sans-serif";
  ctx.textBaseline = "middle";
  ctx.textAlign = "center";
  for (let i = 0; i < rectangles.length; i++) {
    const cr = rect2canvas(rectangles[i]);
    ctx.fillText("" + i, cr.x + cr.w / 2, cr.y + cr.h / 2);
    ctx.strokeText("" + i, cr.x + cr.w / 2, cr.y + cr.h / 2);
  }
}

function updateInput() {
  const searchParams = new URLSearchParams(window.location.search);
  searchParams.set(
    "r",
    rectangles
      .flatMap((r) => [r.x, r.y, r.w, r.h])
      .map((n) => Math.round(n * 1000))
      .join(".")
  );
  const newRelativePathQuery =
    window.location.pathname + "?" + searchParams.toString();
  history.pushState(null, "", newRelativePathQuery);
}

function mousePos(e) {
  const rect = canvas.getBoundingClientRect();
  return {
    x: (e.clientX - rect.left) / (rect.right - rect.left),
    y: (e.clientY - rect.top) / (rect.bottom - rect.top),
  };
}

canvas.addEventListener("mousedown", (e) => {
  if (!imageLoaded || e.button !== 0) {
    return;
  }
  lastPos = mousePos(e);
  const index = rectIndex(lastPos.x, lastPos.y);
  if (index >= 0) {
    movingIndex = index;
    // move to top
  } else {
    isDrawing = true;
    rectangles.push({
      x: lastPos.x,
      y: lastPos.y,
      w: 0,
      h: 0,
    });
  }
});

canvas.addEventListener("mouseup", (e) => {
  if (!imageLoaded || e.button !== 0) {
    return;
  }
  x = 0;
  y = 0;
  isDrawing = false;
  movingIndex = -1;
  lastPos = null;
  updateInput();
});

canvas.addEventListener("mousemove", (e) => {
  let last = rectangles[rectangles.length - 1];
  if (!last) {
    return;
  }
  const { x, y } = mousePos(e);
  if (isDrawing) {
    last.w = x - last.x;
    last.h = y - last.y;
  } else if (movingIndex !== -1) {
    rectangles[movingIndex].x += x - lastPos.x;
    rectangles[movingIndex].y += y - lastPos.y;
  }
  lastPos = { x, y };
  draw();
});

window.addEventListener("keydown", (e) => {
  if (e.key === "Backspace" || e.key === "Delete") {
    rectangles.pop();
    draw();
    updateInput();
  }
});

function rectIndex(x, y) {
  for (let i = rectangles.length - 1; i >= 0; i--) {
    const r = rectangles[i];
    const nr = {
      x: Math.min(r.x, r.x + r.w),
      y: Math.min(r.y, r.y + r.h),
      w: Math.abs(r.w),
      h: Math.abs(r.h),
    };
    if (x >= nr.x && x <= nr.x + nr.w && y >= nr.y && y <= nr.y + nr.h) {
      return i;
    }
  }
  return -1;
}

function uriStringToRects() {
  const searchParams = new URLSearchParams(window.location.search);
  const uriParam = searchParams.get("r");
  if (uriParam === undefined || uriParam === null || uriParam === "") {
    return;
  }
  rectangles = chunks(
    uriParam.split(".").map((i) => i / 1000),
    4
  ).map((r) => ({ x: r[0], y: r[1], w: r[2], h: r[3] }));
}

file.addEventListener(
  "change",
  (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      img.src = e.target.result;
      img.onload = initImage;
      updateInput();
    };
    reader.readAsDataURL(file);
  },
  false
);

uriStringToRects();

window.onresize = initImage;

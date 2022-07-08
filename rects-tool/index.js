const img = new Image();

const container = document.getElementById("container");
const canvas = document.getElementById("canvas");
const file = document.getElementById("file");
const input = document.getElementById("input");
const button = document.getElementById("copy");
const ctx = canvas.getContext("2d");

let x = 0;
let y = 0;
let isDrawing = false;
let isMoving = false;
let rectangles = [];
let lastPos = null;

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
  draw();
}

function chunks(bigarray, size) {
  const arrayOfArrays = [];
  for (let i = 0; i < bigarray.length; i+=size) {
    arrayOfArrays.push(bigarray.slice(i,i+size));
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
  for (r of rectangles) {
    ctx.fillStyle = "white";
    const cr = rect2canvas(r);
    ctx.fillRect(cr.x, cr.y, cr.w, cr.h);
  }
}

function round(n) {
  const accuracy = 1000;
  return Math.round(n * accuracy) / accuracy;
}

function mapRectangles(rectangles) {
  return rectangles.map(normalizeRectangle).map((r) => {
    return {
      min: {
        x: round(r.x),
        y: round(r.y),
      },
      max: {
        x: round(r.x + r.w),
        y: round(r.y + r.h),
      },
    };
  });
}

function updateInput() {
  rectsToUriString();
  input.value = JSON.stringify(mapRectangles(rectangles));
}

function mousePos(e) {
  const rect = canvas.getBoundingClientRect();
  return {
    x: (e.clientX - rect.left) / (rect.right - rect.left),
    y: (e.clientY - rect.top) / (rect.bottom - rect.top),
  };
}

canvas.addEventListener("mousedown", (e) => {
  if (e.button !== 0) {
    return;
  }
  lastPos = mousePos(e);
  const index = rectIndex(lastPos.x, lastPos.y);
  let r;
  if (index >= 0) {
    isMoving = true;
    // move to top
    r = rectangles[index];
    rectangles.splice(index, 1);
    rectangles.push(r);
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
  if (e.button !== 0) {
    return;
  }
  x = 0;
  y = 0;
  isDrawing = false;
  isMoving = false;
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
  } else if (isMoving) {
    last.x += x - lastPos.x;
    last.y += y - lastPos.y;
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

function updateClipboard(newClip) {
  navigator.clipboard.writeText(newClip).then(
    () => {
      button.innerHTML = "Copied!";
    },
    () => {
      button.innerHTML = "Failed to copy!";
    }
  );
  setTimeout(() => {
    button.innerHTML = "Copy";
  }, 2000);
}

function normalizeRectangle(r) {
  return {
    x: Math.min(r.x, r.x + r.w),
    y: Math.min(r.y, r.y + r.h),
    w: Math.abs(r.w),
    h: Math.abs(r.h),
  };
}

function rectIndex(x, y) {
  for (let i = rectangles.length - 1; i >= 0; i--) {
    const r = rectangles[i];
    const nr = normalizeRectangle(r);
    if (x >= nr.x && x <= nr.x + nr.w && y >= nr.y && y <= nr.y + nr.h) {
      return i;
    }
  }
  return -1;
}

button.addEventListener("click", () => {
  const rectsJson = JSON.stringify(mapRectangles(rectangles));
  updateClipboard(rectsJson);
});

function rectsToUriString() {
  if ('URLSearchParams' in window) {
    const searchParams = new URLSearchParams(window.location.search)
    const uriParam = rectangles.map((r) => [r.x, r.y, r.w, r.h].map(round).join("-")).join("-");
    searchParams.set("rects", uriParam);
    const newRelativePathQuery = window.location.pathname + '?' + searchParams.toString();
    history.pushState(null, '', newRelativePathQuery);
  }
}

function uriStringToRects() {
  if ('URLSearchParams' in window) {
    const searchParams = new URLSearchParams(window.location.search);
    const uriParam = searchParams.get("rects");
    if (uriParam === undefined || uriParam === null || uriParam === "") {
      return
    }
    console.log(uriParam)
    rectangles = chunks(uriParam.split('-').map((i) => +i), 4)
        .map((r) => ({x: r[0], y: r[1], w: r[2], h: r[3]}));
    console.log(rectangles);
    updateInput();
  }
}

file.addEventListener(
  "change",
  (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      img.src = e.target.result;
      rectangles = [];
      img.onload = initImage;
      uriStringToRects();
      updateInput();
    };
    result.classList.remove("invisible");
    reader.readAsDataURL(file);
  },
  false
);

input.addEventListener("keyup", (e) => {
  if (e.key === "Enter") {
    try {
      rectangles = JSON.parse(input.value);
      rectangles = rectangles.map((r) => {
        return {
          x: r.min.x,
          y: r.min.y,
          w: r.max.x - r.min.x,
          h: r.max.y - r.min.y,
        };
      });
      draw();
    } catch (e) {
      console.log(e);
    }
  }
});

window.onresize = initImage;
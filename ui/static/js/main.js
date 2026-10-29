// Static demo data. In the real app these come from the Go backend.
const NOTES = {
  K7M2QX: {title:"دکمه‌ی تلاش دوباره در Compose", lang:"Kotlin", age:"۲ دقیقه پیش",
    body:'@Composable\nfun RetryButton(onRetry: () -> Unit) {\n    Button(onClick = onRetry) {\n        Text("تلاش دوباره")\n    }\n}'},
  GOLD42: {title:"هندلر ساده‌ی HTTP در Go", lang:"Go", age:"۱ ساعت پیش",
    body:'func hello(w http.ResponseWriter, r *http.Request) {\n\tfmt.Fprintln(w, "hello")\n}\n\nfunc main() {\n\thttp.HandleFunc("/", hello)\n\thttp.ListenAndServe(":8080", nil)\n}'},
  R3AD9E: {title:"فهرست کتاب‌های آخر هفته", lang:"متن ساده", age:"دیروز",
    body:"۱. Designing Data-Intensive Applications\n۲. The Go Programming Language\n۳. A Philosophy of Software Design"}
};

const ALPHABET = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"; // 32 chars, same as the Go version
const CODE_RE = /^[A-HJ-NP-Z2-9]{6}$/;
const $ = (s) => document.querySelector(s);

const makeCode = () =>
  Array.from(crypto.getRandomValues(new Uint8Array(6)), (b) => ALPHABET[b & 31]).join("");

// Notes created in this browser tab (demo only).
const loadNotes = () => { try { return JSON.parse(sessionStorage.getItem("notes") || "{}"); } catch { return {}; } };
const saveNote = (code, note) => { try { const s = loadNotes(); s[code] = note; sessionStorage.setItem("notes", JSON.stringify(s)); } catch {} };
const findNote = (code) => loadNotes()[code] || NOTES[code];

function copyText(text, btn) {
  const old = btn.textContent;
  (navigator.clipboard ? navigator.clipboard.writeText(text) : Promise.reject())
    .catch(() => {})
    .finally(() => { btn.textContent = "کپی شد"; setTimeout(() => (btn.textContent = old), 1400); });
}

/* ---------- home: open-by-code dialog ---------- */
const dlg = $("#codeDialog");
if (dlg) {
  const input = $("#codeIn"), err = $("#codeErr");
  $("#openBtn").addEventListener("click", () => { err.textContent = ""; input.value = ""; dlg.showModal(); input.focus(); });
  $("#cancelBtn").addEventListener("click", () => dlg.close());
  dlg.addEventListener("click", (e) => { if (e.target === dlg) dlg.close(); });
  $("#codeForm").addEventListener("submit", (e) => {
    e.preventDefault();
    const code = input.value.trim().toUpperCase();
    if (!CODE_RE.test(code)) { err.textContent = "کد باید ۶ حرف یا عدد باشد."; input.focus(); return; }
    location.href = "read_paste.html?code=" + code;
  });
}

/* ---------- create page ---------- */
const form = $("#createForm");
if (form) {
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    const body = $("#body").value;
    if (!body.trim()) { $("#body").focus(); return; }
    const code = makeCode();
    saveNote(code, { title: $("#title").value.trim() || "بدون عنوان", lang: $("#lang").value, age: "همین الان", body });
    $("#newCode").textContent = code;
    $("#viewLink").href = "read_paste.html?code=" + code;
    $("#result").classList.remove("hidden");
    $("#result").scrollIntoView({ behavior: "smooth", block: "center" });
  });
  $("#copyCode").addEventListener("click", (e) => copyText($("#newCode").textContent, e.currentTarget));
}

/* ---------- read page ---------- */
const view = $("#note");
if (view) {
  const code = (new URLSearchParams(location.search).get("code") || "").toUpperCase();
  const note = findNote(code);
  if (note) {
    $("#vTitle").textContent = note.title;
    $("#vCode").textContent = code;
    $("#vLang").textContent = note.lang;
    $("#vAge").textContent = note.age;
    $("#vBody").textContent = note.body;
    view.classList.remove("hidden");
    $("#copyNote").addEventListener("click", (e) => copyText(note.body, e.currentTarget));
  } else {
    $("#notFound").classList.remove("hidden");
  }
}

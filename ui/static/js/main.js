const $ = (s) => document.querySelector(s);

function copyText(text, btn) {
  const old = btn.textContent;
  (navigator.clipboard ? navigator.clipboard.writeText(text) : Promise.reject())
      .catch(() => {})
      .finally(() => { btn.textContent = "کپی شد"; setTimeout(() => (btn.textContent = old), 1400); });
}

/* home: open-by-code dialog */
const dlg = $("#codeDialog");
if (dlg) {
  $("#openBtn").addEventListener("click", () => { $("#codeIn").value = ""; dlg.showModal(); $("#codeIn").focus(); });
  $("#cancelBtn").addEventListener("click", () => dlg.close());
  dlg.addEventListener("click", (e) => { if (e.target === dlg) dlg.close(); });
}

/* view page: copy button */
const copyNote = $("#copyNote");
if (copyNote) {
  copyNote.addEventListener("click", (e) => copyText($("#vBody").textContent, e.currentTarget));
}
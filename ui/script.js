const dropZone = document.getElementById('dropZone');
const fileInput = document.getElementById('fileInput');
const uploadBtn = document.getElementById('uploadBtn');
const uploadList = document.getElementById('uploadList');
const toggleDark = document.getElementById('toggleDark');

let selectedFiles = [];
const MAX_FILE_SIZE_MB = 100; // Max 100MB per file

dropZone.addEventListener('click', () => fileInput.click());

dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropZone.classList.add('border-blue-400');
});

dropZone.addEventListener('dragleave', (e) => {
    e.preventDefault();
    dropZone.classList.remove('border-blue-400');
});

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('border-blue-400');
    selectedFiles = Array.from(e.dataTransfer.files);
    updateDropZoneText();
    uploadBtn.disabled = selectedFiles.length === 0;
});

fileInput.addEventListener('change', (e) => {
    selectedFiles = Array.from(e.target.files);
    updateDropZoneText();
    uploadBtn.disabled = selectedFiles.length === 0;
});

function updateDropZoneText() {
    if (selectedFiles.length > 0) {
        dropZone.querySelector('p').innerText = `${selectedFiles.length} file(s) selected`;
    } else {
        dropZone.querySelector('p').innerText = 'Drag & Drop files here or click to select';
    }
}

uploadBtn.addEventListener('click', async () => {
    if (selectedFiles.length === 0) return;

    uploadBtn.disabled = true;
    uploadBtn.innerText = "Uploading...";

    uploadList.innerHTML = ''; // Clear previous uploads
    for (const file of selectedFiles) {
        if (file.size > MAX_FILE_SIZE_MB * 1024 * 1024) {
            // File too big
            const errorEntry = document.createElement('div');
            errorEntry.className = "bg-red-100 dark:bg-red-800 p-4 rounded shadow text-red-700 dark:text-red-200 mb-2";
            errorEntry.innerText = `❌ ${file.name} is larger than ${MAX_FILE_SIZE_MB}MB. Skipped.`;
            uploadList.appendChild(errorEntry);
            continue;
        }
        await uploadFile(file);
    }

    uploadBtn.innerText = "Upload Files";
    selectedFiles = [];
    updateDropZoneText();
});

async function uploadFile(file) {
    // Create UI entry
    const entry = document.createElement('div');
    entry.className = "bg-gray-100 dark:bg-gray-700 p-4 rounded shadow relative";
    entry.innerHTML = `
    <div class="font-semibold text-gray-800 dark:text-gray-100 mb-2">${file.name}</div>
    <div class="w-full bg-gray-300 dark:bg-gray-600 rounded-full h-2.5 mb-2">
      <div class="bg-blue-600 h-2.5 rounded-full" style="width: 0%" id="progress-${file.name}"></div>
    </div>
    <div class="text-sm text-gray-600 dark:text-gray-400" id="status-${file.name}">Uploading...</div>
  `;
    uploadList.appendChild(entry);

    const progressBar = document.getElementById(`progress-${file.name}`);
    const statusText = document.getElementById(`status-${file.name}`);

    const formData = new FormData();
    formData.append('file', file);

    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/upload', true);

        xhr.upload.addEventListener('progress', (e) => {
            if (e.lengthComputable) {
                const percent = (e.loaded / e.total) * 100;
                progressBar.style.width = `${percent}%`;
            }
        });

        xhr.onreadystatechange = () => {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                const data = JSON.parse(xhr.responseText);
                if (xhr.status === 200) {
                    statusText.innerHTML = `
            <div class="flex items-center gap-2">
              <span>✅ Token:</span> 
              <span class="font-mono" id="token-${file.name}">${data.token}</span> 
              <button class="text-blue-500 hover:text-blue-700 text-xs" onclick="copyToken('${file.name}')">📋 Copy</button>
            </div>
            <div id="countdown-${file.name}" class="text-xs mt-1"></div>
          `;
                    startCountdown(data.expires_at, `countdown-${file.name}`);
                    resolve();
                } else {
                    statusText.innerText = "❌ Upload failed";
                    reject();
                }
            }
        };

        xhr.send(formData);
    });
}

function startCountdown(expireTime) {
    const expireAt = new Date(expireTime); // already ISO, browser parses correctly
    function update() {
        const now = new Date(); // browser local time
        const diff = Math.max(0, expireAt.getTime() - now.getTime());
        const minutes = Math.floor(diff / 60000);
        const seconds = Math.floor((diff % 60000) / 1000);
        countdown.innerText = `Expires in ${minutes}m ${seconds}s`;
        if (diff > 0) {
            setTimeout(update, 1000);
        } else {
            countdown.innerText = "Token expired.";
        }
    }
    update();
}


// 🔥 Copy token function
function copyToken(fileName) {
    const tokenEl = document.getElementById(`token-${fileName}`);
    navigator.clipboard.writeText(tokenEl.innerText).then(() => {
        showToast('✅ Token copied!');
    }).catch(() => {
        showToast('❌ Failed to copy!');
    });
}

// 🛠 Show toast function
function showToast(message) {
    const toast = document.getElementById('toast');
    toast.innerText = message;
    toast.classList.remove('hidden');
    setTimeout(() => {
        toast.classList.remove('opacity-0');
    }, 10); // Start fade-in after slight delay

    setTimeout(() => {
        toast.classList.add('opacity-0');
        setTimeout(() => {
            toast.classList.add('hidden');
        }, 500); // Wait fade-out complete then hide
    }, 2000); // Show for 2 sec
}

const downloadForm = document.getElementById('downloadForm');

downloadForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const tokenInput = document.getElementById('tokenInput');
    const token = tokenInput.value.trim();
    if (!token) return;

    window.location.href = `/download/${token}`;
});

toggleDark.addEventListener('click', () => {
    document.documentElement.classList.toggle('dark');
});

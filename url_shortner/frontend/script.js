// --- Backend API ---

const API_URL = 'http://localhost:8080';


// --- Section switching ---

const shortenerBtn =
    document.getElementById('shortenerBtn');

const customAliasBtn =
    document.getElementById('customAliasBtn');

const shortenerSection =
    document.getElementById('shortenerSection');

const customAliasSection =
    document.getElementById('customAliasSection');


function showSection(which) {

    if (which === 'shortener') {

        shortenerSection.style.display = 'block';
        customAliasSection.style.display = 'none';

    } else {

        shortenerSection.style.display = 'none';
        customAliasSection.style.display = 'block';

    }
}


shortenerBtn.addEventListener('click', () => {

    showSection('shortener');

});


customAliasBtn.addEventListener('click', () => {

    showSection('customAlias');

});


// --- URL validation ---

function isValidURL(value) {

    try {

        const url = new URL(value);

        return (
            (url.protocol === 'http:' ||
             url.protocol === 'https:') &&
            url.hostname !== ''
        );

    } catch {

        return false;

    }

}


// --- Shortener form ---

const shortenerForm =
    document.getElementById('shortenerForm');

const originalUrlInput =
    document.getElementById('originalUrl');

const shortenButton =
    document.getElementById('shortenButton');

const shortenerResult =
    document.getElementById('shortenerResult');

const customPrompt =
    document.getElementById('customPrompt');

const pickAliasBtn =
    document.getElementById('pickAliasBtn');


shortenerForm.addEventListener('submit', async (e) => {

    e.preventDefault();


    const url =
        originalUrlInput.value.trim();


    // Clear previous result

    shortenerResult.innerHTML = '';

    customPrompt.style.display = 'none';


    // Empty URL

    if (!url) {

        showError(
            shortenerResult,
            'Please enter the URL'
        );

        return;
    }


    // Invalid URL

    if (!isValidURL(url)) {

        showError(
            shortenerResult,
            '⚠️ Please enter a valid HTTP or HTTPS URL'
        );

        return;
    }


    // Loading state

    shortenButton.disabled = true;
    shortenButton.textContent = 'Creating...';


    try {

        const response = await fetch(
            `${API_URL}/shorten`,
            {
                method: 'POST',

                headers: {
                    'Content-Type': 'application/json'
                },

                body: JSON.stringify({
                    original_url: url,
                    alias: ''
                })
            }
        );


        const data =
            await readResponse(response);


        // Backend error

        if (!response.ok) {

            showError(
                shortenerResult,
                data.message
            );

            return;
        }


        // Success

        const shortUrl =
            data.short_url;


        shortenerResult.innerHTML = `

            <div class="success-message">
                URL shortened successfully! 🎉
            </div>

            <div class="result-link-row">

                <div class="result-link">

                    <a
                        href="${shortUrl}"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        ${shortUrl}
                    </a>

                </div>

                <button
                    type="button"
                    class="copy-btn"
                    data-copy="${shortUrl}"
                >
                    Copy
                </button>

            </div>

        `;


        customPrompt.style.display = 'block';


    } catch (error) {

        showError(
            shortenerResult,
            '⚠️ Unable to connect to the server.'
        );

    } finally {

        shortenButton.disabled = false;
        shortenButton.textContent = 'Shorten URL';

    }

});


// --- Pick custom alias ---

pickAliasBtn.addEventListener('click', () => {

    const url =
        originalUrlInput.value.trim();


    showSection('customAlias');


    document
        .getElementById('customOriginalUrl')
        .value = url;


    document
        .getElementById('customOriginalUrl')
        .focus();

});


// --- Custom alias form ---

const customAliasForm =
    document.getElementById('customAliasForm');

const customOriginalUrlInput =
    document.getElementById('customOriginalUrl');

const aliasInput =
    document.getElementById('alias');

const customAliasButton =
    document.getElementById('customAliasButton');

const customAliasResult =
    document.getElementById('customAliasResult');


customAliasForm.addEventListener('submit', async (e) => {

    e.preventDefault();


    const url =
        customOriginalUrlInput.value.trim();

    const alias =
        aliasInput.value.trim();


    // Clear previous result

    customAliasResult.innerHTML = '';


    // Both empty

    if (!url && !alias) {

        showError(
            customAliasResult,
            '⚠️ Please enter the URL and alias name you want'
        );

        return;
    }


    // Empty URL

    if (!url) {

        showError(
            customAliasResult,
            '⚠️ Please enter the URL'
        );

        return;
    }


    // Empty alias

    if (!alias) {

        showError(
            customAliasResult,
            '⚠️ Please enter an alias'
        );

        return;
    }


    // Alias length

    if (alias.length > 100) {

        showError(
            customAliasResult,
            '⚠️ Alias must be 100 characters or less'
        );

        return;
    }


    // Invalid URL

    if (!isValidURL(url)) {

        showError(
            customAliasResult,
            '⚠️ Please enter a valid HTTP or HTTPS URL'
        );

        return;
    }


    // Invalid alias characters

    if (!/^[A-Za-z0-9_-]+$/.test(alias)) {

        showError(
            customAliasResult,
            '⚠️ Alias can only contain letters, numbers, _ and -'
        );

        return;
    }


    // Loading state

    customAliasButton.disabled = true;
    customAliasButton.textContent = 'Creating...';


    try {

        const response = await fetch(
            `${API_URL}/shorten`,
            {
                method: 'POST',

                headers: {
                    'Content-Type': 'application/json'
                },

                body: JSON.stringify({
                    original_url: url,
                    alias: alias
                })
            }
        );


        const data =
            await readResponse(response);


        // Backend error

        if (!response.ok) {

            showError(
                customAliasResult,
                data.message
            );

            return;
        }


        // Success

        const shortUrl =
            `${API_URL}/${data.custom_alias}`;


        customAliasResult.innerHTML = `

            <div class="success-message">
                Custom alias created successfully! 🎉
            </div>

            <div class="result-link-row">

                <div class="result-link">

                    <a
                        href="${shortUrl}"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        ${shortUrl}
                    </a>

                </div>

                <button
                    type="button"
                    class="copy-btn"
                    data-copy="${shortUrl}"
                >
                    Copy
                </button>

            </div>

        `;


    } catch (error) {

        showError(
            customAliasResult,
            '⚠️ Unable to connect to the server.'
        );

    } finally {

        customAliasButton.disabled = false;
        customAliasButton.textContent =
            'Create Short URL';

    }

});


// --- Read backend response safely ---

async function readResponse(response) {

    const responseText =
        await response.text();


    if (!responseText) {

        return {};

    }


    try {

        return JSON.parse(responseText);

    } catch {

        return {
            message: responseText
        };

    }

}


// --- Show error ---

function showError(element, message) {

    element.innerHTML = `

        <div class="error-message">
            ${message || '⚠️ Something went wrong.'}
        </div>

    `;

}


// --- Copy button ---

const copyToast =
    document.getElementById('copyToast');


document.addEventListener('click', async (e) => {

    if (!e.target.matches('.copy-btn')) {

        return;
    }


    const text =
        e.target.getAttribute('data-copy');


    try {

        await navigator.clipboard.writeText(text);

        showCopiedMessage();

    } catch {

        showError(
            e.target.parentElement.parentElement,
            '⚠️ Unable to copy the URL.'
        );

    }

});


function showCopiedMessage() {

    copyToast.classList.add('show');


    setTimeout(() => {

        copyToast.classList.remove('show');

    }, 1500);

}


// --- Default view ---

showSection('shortener');
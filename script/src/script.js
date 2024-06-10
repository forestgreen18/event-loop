const gfButton = document.querySelector(".gf-script");
const dmButton = document.querySelector(".dm-script");

// Base URL for HTTP requests
const baseURL = "http://127.0.0.1:17000/";

// Selectors remain unchanged
const form = document.querySelector("form");
const textarea = document.querySelector("textarea");

// Function to build the full URL with query parameters
function buildURL(params) {
  const queryParams = params
    .map((param) => encodeURIComponent(param))
    .join(",");
  return `${baseURL}?cmd=${queryParams}`;
}

// Refactored HTTP request function using async/await
async function sendHTTPRequest(url) {
  try {
    const response = await fetch(url);
    if (!response.ok)
      throw new Error(`Request failed with status ${response.status}`);
    return await response.text();
  } catch (error) {
    console.error("HTTP Request Failed:", error);
  }
}

// Event listener for form submission
form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const params = textarea.value.trim().split("\n");
  const url = buildURL(params);
  try {
    const response = await sendHTTPRequest(url);
    console.log(response);
  } catch (error) {
    console.error(error);
  }
});

// Event listener for 'Green Fill' button
gfButton.addEventListener("click", async () => {
  const url = buildURL(["green", "bgrect 0.05 0.05 0.95 0.95", "update"]);
  try {
    const response = await sendHTTPRequest(url);
    console.log(response);
  } catch (error) {
    console.error(error);
  }
});

// Event listener for 'Draw and Move' button
dmButton.addEventListener("click", () => {
  const urlToDraw = buildURL(["white", "figure 0.1 0.1", "update"]);
  sendHTTPRequest(urlToDraw)
    .then((response) => console.log(response))
    .catch((error) => console.error(error));

  for (let i = 0; i < 9; i++) {
    setTimeout(() => {
      console.log(`Request: ${i + 1}`);
      const urlToMove = buildURL(["move 0.1 0.1", "update"]);
      sendHTTPRequest(urlToMove)
        .then((response) => console.log(response))
        .catch((error) => console.error(error));
    }, (i + 1) * 1000);
  }
});

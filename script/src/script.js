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

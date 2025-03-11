const fs = require("fs");

const url = "http://localhost:8080/generate";
const payload = {
  language: "Python",
};

fetch(url, {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify(payload),
})
  .then((response) => response.json())
  .then((data) => {
    fs.writeFileSync("python-course.txt", JSON.stringify(data, null, 2));
    console.log("✅ Course saved to python-course.txt");
  })
  .catch((err) => console.error("❌ Error:", err));

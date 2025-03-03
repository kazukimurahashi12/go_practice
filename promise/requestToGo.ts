import axios from "axios";

// サーバサイドエンドポイント
// 画像アップロードリクエスト
app.post("/upload", async (req, res) => {
    const { imageUrl } = req.body;

    await requestImageProcessing(imageUrl);

    res.json({ message: "Image processing started." });
  });

  // 画像処理
async function requestImageProcessing(imageUrl: string) {
    // Go側へリクエスト実行
  const response = await axios.post("http://localhost:8080/process", { imageUrl });

  console.log("Processing started:", response.data);
}



import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import okhttp3.*
import java.io.IOException

class MainActivity : AppCompatActivity() {
    private val client = OkHttpClient()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        fetchCourse()
    }

    private fun fetchCourse() {
        val url = "http://localhost:8080/generate"
        val json = """
            {
                "language": "Python"
            }
        """.trimIndent()

        val requestBody = RequestBody.create(MediaType.get("application/json; charset=utf-8"), json)
        val request = Request.Builder()
            .url(url)
            .post(requestBody)
            .build()

        client.newCall(request).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                println("❌ Error: $e")
            }

            override fun onResponse(call: Call, response: Response) {
                val responseBody = response.body?.string()
                println("✅ Course: $responseBody")
            }
        })
    }
}

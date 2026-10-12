#include "http.h"

#include <libsoup/soup.h>

#include <cstdlib>
#include <stdexcept>

std::string base_url() {
  const char* env = std::getenv("AI_PLANNER_API");
  return env ? std::string(env) : std::string("http://planner.lucia.friedmicro-lab.org");
}

std::string http_json(const std::string& method, const std::string& path,
                      const std::string& body) {
  g_autoptr(SoupSession) session = soup_session_new();
  g_autoptr(SoupMessage) msg =
      soup_message_new(method.c_str(), (base_url() + path).c_str());
  if (msg == nullptr) throw std::runtime_error("soup_message_new failed: " + path);

  if (!body.empty()) {
    g_autoptr(GBytes) bytes = g_bytes_new(body.data(), body.size());
    soup_message_set_request_body_from_bytes(msg, "application/json", bytes);
  }

  g_autoptr(GError) err = nullptr;
  g_autoptr(GBytes) out = soup_session_send_and_read(session, msg, nullptr, &err);
  if (err) throw std::runtime_error("HTTP " + method + " " + path + ": " + err->message);

  gsize status = (gsize)soup_message_get_status(msg);
  gsize len = 0;
  const char* data = static_cast<const char*>(g_bytes_get_data(out, &len));
  std::string text(data ? data : "", len);
  if (status < 200 || status >= 300) {
    throw std::runtime_error("HTTP " + std::to_string(status) + " " + method + " " + path + ": " + text);
  }
  return text;
}

JsonNode* parse_json(const std::string& text) {
  if (text.empty()) return nullptr;
  g_autoptr(GError) err = nullptr;
  g_autoptr(JsonParser) parser = json_parser_new();
  if (!json_parser_load_from_data(parser, text.c_str(), text.size(), &err)) {
    throw std::runtime_error("JSON parse failed: " + std::string(err->message));
  }
  return json_node_ref(json_parser_get_root(parser));
}

double jnum(JsonNode* n) { return n && json_node_get_value_type(n) == G_TYPE_DOUBLE ? json_node_get_double(n) : 0; }
int jint(JsonNode* n) { return n && json_node_get_value_type(n) == G_TYPE_INT64 ? (int)json_node_get_int(n) : 0; }

std::string jstr(JsonNode* n) {
  if (!n) return "";
  GType t = json_node_get_value_type(n);
  if (t == G_TYPE_STRING) {
    const char* s = json_node_get_string(n);
    return s ? s : "";
  }
  if (t == G_TYPE_INT64) return std::to_string(json_node_get_int(n));
  if (t == G_TYPE_DOUBLE) return std::to_string(json_node_get_double(n));
  return "";
}

std::vector<int> jday_array(JsonNode* node) {
  std::vector<int> out;
  if (!node || !JSON_NODE_HOLDS_ARRAY(node)) return out;
  JsonArray* arr = json_node_get_array(node);
  guint len = json_array_get_length(arr);
  for (guint i = 0; i < len; ++i) out.push_back((int)json_array_get_int_element(arr, i));
  return out;
}

JsonObject* as_object(JsonNode* root) {
  if (!root || !JSON_NODE_HOLDS_OBJECT(root)) return nullptr;
  return json_node_get_object(root);
}
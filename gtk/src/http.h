#ifndef PLANNER_HTTP_H
#define PLANNER_HTTP_H

#include <json-glib/json-glib.h>

#include <string>
#include <vector>

std::string base_url();
std::string http_json(const std::string& method, const std::string& path,
                      const std::string& body = "");
JsonNode* parse_json(const std::string& text);

double jnum(JsonNode* n);
int jint(JsonNode* n);
std::string jstr(JsonNode* n);
std::vector<int> jday_array(JsonNode* node);
JsonObject* as_object(JsonNode* root);

#endif

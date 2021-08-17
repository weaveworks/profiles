/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Modifications:
  Copyright 2021 Weaveworks
*/

import React, { useState, useEffect } from "react";

const Schema = ({jsonFile}) => {
  const [data, setData] = useState();

  const valueEntry = (definition) => {
    let value = definition.default;
    let valueClass = "value";
    let tooltip = "default";
    const isEnum =
      (definition.enum && definition.enum.length > 0) || definition.const;
    if (definition.default == null && isEnum) {
      value = definition.const || definition.example || definition.enum[0];
      valueClass = "example";
      tooltip = "example";
      if (definition.const || definition.enum.length === 1) {
        valueClass = "const";
        tooltip = "required value";
      }
    } else if (definition.examples && definition.examples.length > 0) {
      value = definition.examples[0];
      valueClass = "example";
      tooltip = "example";
    }
    return [value, valueClass, tooltip];
  };

  const anchor = (path, label, firstOfListType) => {
    let listPrefix;
    if (firstOfListType !== undefined) {
      listPrefix = (
        <span style={{ visibility: firstOfListType ? "visible" : "hidden" }}> - </span>
      );
    }
    const href = "#" + path;
    const anchor = (
      <a className="key" href={href}>
        {label}
      </a>
    );
    return (
      <>
        {listPrefix}
        {anchor}
      </>
    );
  };

  const offset = (ident) => `${ident * 2}ex`;

  useEffect(() => {
    let isMounted = true;
    fetch(jsonFile)
      .then((resp) => resp.json())
      .then(json => {
        if (isMounted) setData(json);
      })
      .catch((error) => console.log(error));
    return () => { isMounted = false };
  }, [jsonFile]);

  if (!data) {
    return null;
  }

  const getTableRows = (definitions, parentDefinition, ref, ident, parent) => {
    const name = ref.replace("#/definitions/", "");

    const allProperties = [];
    const seen = {};

    for (const key of definitions[name].preferredOrder || []) {
      allProperties.push([key, definitions[name].properties[key]]);
      seen[key] = true;
    }

    let index = -1;

    let result;

    for (let [key, definition] of allProperties) {
      const path = parent.length == 0 ? key : `${parent}-${key}`;
      index++;

      // Key
      const required =
        definitions[name].required && definitions[name].required.includes(key);
      let keyClass = required ? "key required" : "key";

      // Value
      let [value, valueClass, tooltip] = valueEntry(definition);

      // Description
      let desc = definition["x-intellij-html-description"] || "";

      let firstOfListType = undefined;
      if (parentDefinition && parentDefinition.type === "array") {
        firstOfListType = index === 0;
      }

      // Value Cell
      const valueCell = value && (
        <span title={tooltip} className={valueClass}>
          {value}
        </span>
      );

      const keyCell = (
        <td>
          <div className="anchor" id={path}></div>
          <span
            title={required ? "Required key" : ""}
            className={keyClass}
            style={{ marginLeft: offset(ident) }}
          >
            {anchor(path, key, firstOfListType)}:&nbsp;
          </span>
          {valueCell}
        </td>
      );

      // Whether our field has sub fields
      let ref;
      // This definition references another definition directly
      if (definition.$ref) {
        ref = definition.$ref;
        // This definition is an array
      } else if (definition.items && definition.items.$ref) {
        ref = definition.items.$ref;
      }

      if (definition.$ref) {
        // Check if the referenced description is a final one
        const refName = definition.$ref.replace("#/definitions/", "");
        const refDef = definitions[refName];
        let type = "";

        if (refDef.type === "object") {
          if (!refDef.properties && !refDef.anyOf) {
            type = "object";
          }
        } else {
          type = refDef.type;
        }
        if (desc === "") {
          desc = refDef["x-intellij-html-description"] || "";
        }
        result = (
          <>
            {result}
            <tr className="top">
              {keyCell}
              <td className="comment">{desc + " "}</td>
              <td className="type">{type}</td>
            </tr>
          </>
        );
      } else if (definition.items && definition.items.$ref) {
        const refName = definition.items.$ref.replace("#/definitions/", "");
        const refDef = definitions[refName];
        let type = "";
        if (refDef.type === "object") {
          if (!refDef.properties && !refDef.anyOf) {
            type = "object[]";
            value = "{}";
          }
        } else {
          type = `${refDef.type}[]`;
        }
        // If the ref has enum information, show it in the field
        if (desc === "" || (refDef.enum && refDef.enum.length > 0)) {
          desc = [desc, refDef["x-intellij-html-description"]]
            .filter((x) => x)
            .join(" ");
        }
        if (type === "undefined[]") {
          type = "[]";
        }
        result = (
          <>
            {result}
            <tr className="top">
              {keyCell}
              <td className="comment">{desc}</td>
              <td className="type">{type}</td>
            </tr>
          </>
        );
      } else if (definition.type === "array" && value && value !== "[]") {
        // Parse value to json array
        const values = JSON.parse(value);
        const valuesTop = (
          <tr>
            {keyCell}
            <td className="comment" rowspan={1 + values.length}>
              {desc}
            </td>
            <td className="type"></td>
          </tr>
        );
        const valuesInfo = values.map((v) => (
          <tr>
            <td>
              <span className="key" style={{ marginLeft: offset(ident) }}>
                - <span className="valueClass">{v}</span>
              </span>
            </td>
            <td className="comment"></td>
            <td className="type">object</td>
          </tr>
        ));
        result = (
          <>
            {result}
            {valuesTop}
            {valuesInfo}
          </>
        );
      } else if (definition.type === "object" && value && value !== "{}") {
        result = (
          <>
            {result}
            <tr>
              {keyCell}
              <td className="comment">{desc}</td>
              <td className="type">object</td>
            </tr>
          </>
        );
      } else {
        const type =
          definition.type === "array"
            ? `${definition.items.type}[]`
            : definition.type;
        result = (
          <>
            {result}
            <tr>
              {keyCell}
              <td className="comment">{desc}</td>
              <td className="type">{type}</td>
            </tr>
          </>
        );
      }
      if (ref) {
        result = (
          <>
            {result}
            {getTableRows(definitions, definition, ref, ident + 2, path)}
          </>
        );
      }
    }
    return result;
  };

  return (
    <table id="schema">
      <tbody>
        {getTableRows(data.definitions, undefined, data.$ref, 0, "")}
      </tbody>
    </table>
  );
};

export default Schema;

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

  Note: This code is adapted to work with docusaurus specific client 
  and server side rendering. 
  This may NOT work with other documentation frameworks.
*/

import React, {useState, useEffect, useRef} from 'react';
import ExecutionEnvironment from '@docusaurus/ExecutionEnvironment';

let html = () => null;
let render = () => null;
let unsafeHTML = () => null;

if (ExecutionEnvironment.canUseDOM) {
    html = require('lit-html').html
    render = require('lit-html').render
    unsafeHTML = require('lit-html/directives/unsafe-html.js').unsafeHTML
}

function offset(ident) {
    return `${ident * 2}ex`;
}

function valueEntry(definition) {
    let value = definition.default;
    let valueClass = "value";
    let tooltip = "default";
    const isEnum = (definition.enum && definition.enum.length > 0) || definition.const;
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
}

function* template(definitions, parentDefinition, ref, ident, parent) {
    const name = ref.replace("#/definitions/", "");
    const allProperties = [];
    const seen = {};

    for (const key of definitions[name].preferredOrder || []) {
        allProperties.push([key, definitions[name].properties[key]]);
        seen[key] = true;
    }

    let index = -1;
    for (let [key, definition] of allProperties) {
        const path = parent.length == 0 ? key : `${parent}-${key}`;
        index++;

        // Key
        const required =
            definitions[name].required &&
            definitions[name].required.includes(key);
        let keyClass = required ? "key required" : "key";

        // Value
        let [value, valueClass, tooltip] = valueEntry(definition);

        // Description
        let desc = definition["x-intellij-html-description"] || "";

        let firstOfListType = undefined;
        if (parentDefinition && parentDefinition.type === "array") {
            firstOfListType = index === 0;
        }
        const valueCell = value != null
            ? html`<span title="${tooltip}" class="${valueClass}"
                  >${value}</span
              >`
            : null;

        const keyCell = html`
            <td>
                <div class="anchor" id="${path}"></div>
                <span title="${required ? 'Required key' : ''}" class="${keyClass}" style="margin-left: ${offset(ident)}">
                    ${anchor(path, key, firstOfListType)}:
                </span>
                ${valueCell}
            </td>
        `;

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

            yield html`
                <tr class="top">
                    ${keyCell}
                    <td class="comment">${unsafeHTML(desc)}</td>
                    <td class="type">${type}</td>
                </tr>
            `;
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
                desc = [desc, refDef["x-intellij-html-description"]].filter(x => x).join(" ");
            }

            // skip undefined types
            if (type === "undefined[]") {
                type = "[]";
            }

            yield html`
                <tr class="top">
                    ${keyCell}
                    <td class="comment">${unsafeHTML(desc)}</td>
                    <td class="type">${type}</td>
                </tr>
            `;
        } else if (definition.type === "array" && value && value !== "[]") {
            // Parse value to json array
            const values = JSON.parse(value);

            yield html`
                <tr>
                    ${keyCell}
                    <td class="comment" rowspan="${1 + values.length}">
                        ${unsafeHTML(desc)}
                    </td>
                    <td class="type"></td>
                </tr>
            `;

            for (const v of values) {
                yield html`
                    <tr>
                        <td>
                            <span
                                class="key"
                                style="margin-left: ${offset(ident)}"
                                >- <span class="${valueClass}">${v}</span></span
                            >
                        </td>
                        <td class="comment"></td>
                        <td class="type">object</td>
                    </tr>
                `;
            }
        } else if (definition.type === "object" && value && value !== "{}") {
            yield html`
                <tr>
                    ${keyCell}
                    <td class="comment">${unsafeHTML(desc)}</td>
                    <td class="type">object</td>
                </tr>
            `;
        } else {
            const type =
                definition.type === "array"
                    ? `${definition.items.type}[]`
                    : definition.type;
            yield html`
                <tr>
                    ${keyCell}
                    <td class="comment">${unsafeHTML(desc)}</td>
                    <td class="type">${type}</td>
                </tr>
            `;
        }

        // This definition references another definition
        if (ref) {
            yield html`
                ${template(definitions, definition, ref, ident + 2, path)}
            `;
        }
    }
}

function anchor(path, label, firstOfListType) {
    let listPrefix = "";
    if (firstOfListType !== undefined) {
        listPrefix = html`<span
            style="visibility: ${firstOfListType ? "visible" : "hidden"}"
            >-
        </span>`;
    }
    const anchor = html`<a class="key" href="#${path}">${label}</a>`;
    return html`${listPrefix}${anchor}`;
}

export default function Schema({jsonFile}) {
    const [data, updateData] = useState(); //store fetch result
    const refEl = useRef(null); //ref to underlying dom node
  
    //get the schema from the func params
    useEffect(async () => {
      if (ExecutionEnvironment.canUseDOM) {
        const resp = await fetch(jsonFile);
        const json = await resp.json()
        updateData(json);
      }
    }, []);
  
    //render the custom html inside the dom node
    useEffect(() => {
      if (data && refEl.current && ExecutionEnvironment.canUseDOM) {
          render(
              html` ${template(data.definitions, undefined, data.$ref, 0, "")} `,
              refEl.current
          )
      }
    }), [data, refEl]
  
    //react shell dom element
    return ExecutionEnvironment.canUseDOM ? <table id="schema" ref={refEl}></table> : null;
}

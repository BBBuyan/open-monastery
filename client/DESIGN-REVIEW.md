# Design review: Open Monastery

## Summary

The frontend now centers on people in Mongolia finding prayers and rituals for monks to perform on their behalf. Its identity remains the original monastery emblem, green palette, and serif display typography. The previous study-library interpretation was incorrect and has been replaced across Home, Prayers, About, search, feedback states, and page metadata.

This Vue web app uses Apple's accessibility, hierarchy, typography, color, and writing foundations. Native Apple navigation conventions are not requirements for this web interface. Product facts and research belong in [the product context](../docs/product-context.md); domain definitions belong in [the glossary](../GLOSSARY.md).

## Improvements

- **Books and paginated browsing:** The product owner prefers “Books” as the catalog label. Browsing now loads 20 entries per batch using the server's cursor, automatically appends the next batch near the bottom, and provides a keyboard-accessible load-more/retry action. Search queries the complete catalog through the server rather than filtering only loaded pages. Existing entries remain visible on subsequent-page failures. Seven pagination tests and the production build pass; browser scroll behavior has not been visually verified.

- **High — Product meaning, corrected:** `client/src/App.vue` explains monks' recitations and dedications for people and loved ones. `client/src/components/About.vue` explains requester versus beneficiary, religious intentions, and the current availability of online requests. `writing.md › Getting started`: “Determine your app’s voice.” The audience and task now determine the vocabulary.
- **Medium — Catalog meaning, corrected:** `client/src/components/Readings.vue` presents prayers and rituals with their original Mongolian names and descriptions, explicit “Traditional purpose” labels, and `lang="mn"`. Search explains that it matches Mongolian names and purposes. Statements about spiritual benefit are attributed to Buddhist tradition.
- **Medium — Action clarity, retained:** The primary link says “Find a prayer” and opens the catalog. A notice explains that online requests are not available yet. `buttons.md › Content`: “Ensure that each button clearly communicates its purpose.” No action suggests a request has been submitted, scheduled, or paid.
- **Medium — Feedback, retained:** Loading, retry, empty, and no-match states use prayer terminology and provide relevant next steps. `writing.md › Best practices`: “Provide clear next steps on any blank screens.”

## Accessibility and craft

The earlier accessibility changes remain: visible keyboard focus, ordinary web navigation links with current-page state, skip navigation, controls at least 44px high, relative typography, responsive grids, system dark mode, increased-contrast borders, and reduced-motion handling. CSS pixels are web units, not native points. Original muted text #6e746c on #f7f4ec measured 4.36:1; its replacement #596354 measures 5.72:1. Dark secondary text #b8c2b3 on #171e19 measures 9.23:1. Ratios are calculated from source hex values.

Judgment: The existing visual identity is appropriate to retain; a new decorative theme would not resolve the underlying product misunderstanding. The hero's single job is explaining prayers performed by monks for people. Compact layouts stack the introduction and contextual note; regular layouts place them alongside each other. The inactive login accessory remains removed. No extra animation or glass is needed.

## Validation and limits

- `npm.cmd run build` passes Vue TypeScript checks and Vite production compilation after the context corrections.
- Changed frontend files and agent documentation are checked for whitespace errors separately from pre-existing server edits.
- Catalog review used `server/monastery_books.json`, with representative prayers for children, livestock, wellbeing, overcoming obstacles, and the deceased. Original catalog text is unchanged.
- Visual browser checks and assistive-technology checks have not been performed. The catalog's live API availability has not been verified in this pass.
- A server contract mismatch blocks enabling a dependable request flow; details are recorded in the product context. Currency and tariff semantics also remain unconfirmed. Existing server edits were left untouched.

Apple references above are the installed skill's local files under `C:/Users/batbu/.agents/skills/apple-design/references/hig/`.

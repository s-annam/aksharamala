# History of Aksharamala

# **Introduction**

Aksharamala is a **transliteration system** that dates back to the **glorious days of dot-matrix printers and MS-DOS**. Originally developed in **1989** as a college project—when transliteration wasn't exactly a hot topic—it allowed typing in Indian languages on personal computers. The founders, fueled by passion (and probably some late-night coding marathons), revisited the idea in the early 2000s, building a **Windows application using system hooks, C++, and MFC**.

After a long dormancy, `Aksharamala` is undergoing a **complete revival in 2025**, rewritten in **Go** with **JSON-based keymaps** for flexibility, extensibility, and integration into modern AI-driven environments."

---

# **Timeline**

* **1989**: Developed as a college project **in Turbo Pascal**.  
* **2000-2003**: Expanded into a commercial product using **C++ and MFC**; active development and product growth.  
* **2003-2006**: Limited activity focused on **bug fixes and customer support**.  
* **2006-2024**: Development **ceased** as the **main founder joined Google**, leading to a **potential conflict of interest** (and, let’s be honest, **not enough bandwidth**). Complete shutdown of activity.  
* **2025**: **Complete rewrite in Go**, introducing modern features like **JSON-based keymaps**. The project is **now open source** under [GNU Affero General Public License v3.0 or later](https://github.com/s-annam/aksharamala/blob/main/LICENSE).

---

# **Detailed History**

## **The College Project (1989)**

Ah, **1989**—the era of **beige computers, floppy disks, and dot-matrix symphonies**. While personal computing was just taking off in **India**, a group of passionate students **(who primarily wanted to type in Telugu)** created a system that could render Indian scripts.

The program was written in **Turbo Pascal**, a language that was all the rage in the pre-Windows world. Since **text mode didn’t support complex scripts**, the program had to **switch to graphics mode to display characters**—an innovative (if slightly clunky) solution at the time. Printing was another adventure altogether, requiring **dot-matrix printers to slowly scream out** the characters one line at a time.

While this project was meant to be a **one-off college assignment**, the name **Aksharamala** was born—and so was the group’s early foray into **transliteration technology**.

---

## **Aksharamala Reincarnation (2000-2002)**

By the early 2000s, **English dominated computing**, leaving non-English, and more especially Indic, users struggling with clunky workarounds. Sensing a **huge gap**, the **founders of Aksharamala** envisioned a **commercial solution** that would be:

* **Intuitive** (because no one wants to learn 50 keyboard shortcuts to type).  
* **Standards-compliant** (**Unicode-based**, unlike many hacky ASCII font solutions at the time).  
* **Application-independent**, allowing input in **any Windows program**.

The first commercial versions of Aksharamala **supported all Indic languages available in Unicode at the time**, including **Hindi, Telugu, Marathi, Tamil, Bengali, Kannada, Malayalam, and even Urdu**—a language often overlooked in many Indic input systems of that era.

Additionally, **Aksharamala actively supported language-to-language conversion**, enabling users to transform text from **Sanskrit to Hindi, Telugu, or other Indic scripts**. This cross-language mapping provided a valuable tool for scholars, researchers, and digital content creators.

---

## **Growth and Challenges (2002-2003)**

As Aksharamala gained traction, development focused on **expanding language support** and improving usability. However, **challenges** quickly emerged:

* **Unicode support on Windows** was still maturing, leading to unexpected quirks.  
* **Legacy fonts** posed issues, making transitions to Unicode less than smooth.  
* **Indian language input methods weren’t standardized**, leading to debates over the best approach.

Despite these challenges, **Aksharamala adapted** through **user feedback and technical innovations**, even exploring **collaborations with research institutions**.

---

## **Sustaining the Product (2003-2006)**

By 2003, Aksharamala had become **stable and usable**, shifting the focus from **new features to maintenance**:

* **Bug fixes and updates** were rolled out.  
* **Customer support** ensured existing users weren’t stranded.  
* **Some exploratory work** on expanding to more Indian languages was done, though no major releases followed.

By **2006**, active development had mostly **wound down**.

---

## **A Long Hiatus (2006-2024)**

And then… **radio silence.**

In **2006**, the **main founder joined Google**, bringing with it a **potential conflict of interest** (and let’s be honest, **not enough free time**). As a result, **Aksharamala development completely ceased**. While some die-hard users **kept using the old versions**, there were **no updates, no releases, and no new features**.

Despite this hiatus, the world moved on. **Unicode adoption improved, mobile devices became the primary computing platforms, and AI-driven text processing became mainstream**. However, transliteration systems—especially those handling **cross-language mappings and reverse conversions**—did not evolve as much as they could have.

---

## **Transition to Open Source (2025-Present)**

After nearly **two decades of dormancy**, Aksharamala is **back**—this time as an **open-source project** ([GitHub Repository](https://github.com/s-annam/aksharamala)). The 2025 revival is not just a **reskin of old code**—it’s a **complete rewrite in Go**, bringing:

* **Portability and efficiency**, with modern software development practices.  
* **JSON-based keymaps**, allowing for **dynamic and flexible transliteration rules**.  
* **Cloud and API-friendly architecture**, making it easy to integrate with modern applications.

But **why now?**

We are in the midst of an **AI-driven revolution** in text processing. Large Language Models (LLMs) can generate, translate, and understand text better than ever before—but **script input and transformation tools still lag behind**.

With **Aksharamala’s revival**, we hope to explore:

* How **reversliteration** (going from native scripts **back to phonetic Roman representations**) can be **useful for AI, speech recognition, and linguistic research**.  
* Whether a **fully Unicode-based, API-driven transliteration system** can find applications in **machine learning, NLP, and web development**.  
* What **new use cases emerge** when transliteration is made truly **modular, extensible, and open-source**.

Simply put: **We don’t have all the answers yet, but we’re excited to push the boundaries of transliteration, AI, and NLP—and we welcome contributors, researchers, and developers to join this journey!**

---

# **Looking Ahead**

The open-source journey of Aksharamala is just beginning. The roadmap includes:

* **Expanding support** for **modern operating systems** and **web platforms**.  
* **Enhancing transliteration** for **more Indian languages**.  
* **Developing a robust API**, so **other applications can seamlessly integrate Aksharamala**.  
* **Experimenting with reversliteration, AI integration, and machine learning applications**.

With an **active developer community**, the hope is that **Aksharamala will continue to grow and adapt**, making Indian language computing **as easy and natural as typing in English**.

---

# **Acknowledgments**

Over the years, **many individuals have contributed** to Aksharamala’s development—offering **ideas, feedback, and technical expertise**. While this document doesn’t list names, their efforts are **deeply appreciated**. Aksharamala is a **testament to collective innovation**, and this new chapter **welcomes contributors** to help shape its future.
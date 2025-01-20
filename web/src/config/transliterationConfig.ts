export interface TransliterationScheme {
  id: string;
  name: string;
  description: string;
  isAdvanced?: boolean;
  helpTips?: string[];
}

export interface Language {
  id: string;
  name: string;
  schemes: TransliterationScheme[];
}

export const languages: Language[] = [
  {
    id: 'hindi',
    name: 'Hindi',
    schemes: [
      {
        id: 'hindi',
        name: 'ITRANS',
        description: 'Convert ITRANS text to Devanagari Unicode',
        helpTips: [
          'Type "namaste" to get नमस्ते',
          'Use capital letters for aspirated consonants (e.g., "kh" for ख)',
        ],
      },
      {
        id: 'rhindi',
        name: 'Unicode',
        description: 'Convert Devanagari Unicode to ITRANS',
        isAdvanced: true,
      },
    ],
  },
  {
    id: 'marathi',
    name: 'Marathi',
    schemes: [
      {
        id: 'marathi',
        name: 'ITRANS',
        description: 'Convert ITRANS text to Devanagari Unicode',
        helpTips: [
          'Type "namaste" to get नमस्ते',
          'Use "L" for ळ (retroflex lateral)',
        ],
      },
    ],
  },
  {
    id: 'telugu',
    name: 'Telugu',
    schemes: [
      {
        id: 'TeluguRts',
        name: 'RTS',
        description: 'Convert RTS text to Telugu Unicode',
        helpTips: [
          'Type "namastE" to get నమస్తే',
          'Use "~" for Telugu-specific characters',
        ],
      },
    ],
  },
];

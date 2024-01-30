export type Workbook = {
  id: number;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
};

// export type Problem = {
//   // id: number;
//   type: string;
//   properties: { [key: string]: string };
// };

export type EnglishSentence = {
  srcLang2: string;
  srcText: string;
  srcAudioContent: string;
  srcAudioLengthMillisecond: number;
  dstLang2: string;
  dstText: string;
  dstcAudioContent: string;
  dstAudioLengthMillisecond: number;
};
export type EnglishSentences = {
  sentences: EnglishSentence[];
};
export type WorkbookWithProblems = {
  id: number;
  name: string;
  description: string;
  englishSentences: EnglishSentences;
};

-- Table des mots (Word)
CREATE TABLE WORD
(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    Kanji         VARCHAR(50),
    Onyomi        VARCHAR(50),
    Kunyomi       VARCHAR(50),
    ImageURL      VARCHAR(255),
    TranslationID VARCHAR(50),
    FOREIGN KEY (TranslationID) REFERENCES LABEL (ID)
);

-- Table des étiquettes/traductions (Label)
CREATE TABLE LABEL
(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    en VARCHAR(255),
    fr VARCHAR(255)
);

-- Table d'association pour les tags de mots
CREATE TABLE WORD_TAG
(
    WordID  UUID,
    LabelID UUID,
    FOREIGN KEY (WordID) REFERENCES WORD (ID),
    FOREIGN KEY (LabelID) REFERENCES LABEL (ID),
    PRIMARY KEY (WordID, LabelID)
);

-- Table d'association pour les niveaux de mots
CREATE TABLE WORD_LEVEL
(
    WordID  UUID,
    LabelID UUID,
    FOREIGN KEY (WordID) REFERENCES WORD (ID),
    FOREIGN KEY (LabelID) REFERENCES LABEL (ID),
    PRIMARY KEY (WordID, LabelID)
);

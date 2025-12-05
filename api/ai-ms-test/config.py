import os
from dotenv import load_dotenv
from typing import TypedDict
from langchain_groq import ChatGroq


class State(TypedDict):
    input: str
    output: str
    sql_querry:str

os.environ["OPENAI_API_KEY"]="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"
load_dotenv()

class Config:
    GROQ_API_KEY=os.getenv("GROQ_API_KEY")
    base_llm=ChatGroq(
        name="base_llm",
        model="llama-3.1-8b-instant",
        temperature=0.2,
        max_tokens=500,
        timeout=None,
        max_retries=5
        )
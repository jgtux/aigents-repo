from langchain.agents.xml.prompt import agent_instructions

from config import State, Config
from langchain_groq import ChatGroq
from langchain.agents import create_react_agent, AgentExecutor
from langchain.prompts import ChatPromptTemplate
from langchain.tools import tool


class AgentLLMCore:
    def __init__(self):
        self.stored_llm={"base_llm":Config.base_llm}
        self.stored_agents={}

    def get_all_llm_names(self):
        aux=[]
        for l in self.stored_llm:
            aux.append(l)
        return aux

    def get_one_llm(self,name:str):
        for l in self.stored_llm:
            if l == name:
                return self.stored_llm[l]

        print("llm not found")
        return None
    def get_all_agents(self):
        aux=[]
        for i in self.stored_agents:
            aux.append(i)
        return aux

    def create_llm(self, name:str, model:str ="llama-3.1-8b-instant",temperature:float =0.2,max_tokens:int =500,timeout=None):
        aux=ChatGroq(name=name,model=model,temperature=temperature,max_tokens=max_tokens,timeout=timeout)
        self.stored_llm[name]=aux
        return aux

    def create_agent(self, llm, tools, prompt_text: str, name: str):
        # por enquanto ignoramos tools; foco no LLM simples
        if isinstance(llm, str):
            aux_llm = self.get_one_llm(llm)
        else:
            aux_llm = llm

        if aux_llm is None:
            raise ValueError(f"LLM '{llm}' não encontrado")

        # monta um prompt simples: suas instruções + input do usuário
        template = f"""{prompt_text}

        Usuário: {{input}}
        """

        final_prompt = ChatPromptTemplate.from_template(template)

        # chain = prompt -> llm
        chain = final_prompt | aux_llm

        # colocar um nome pra log/identificação
        chain.name = name
        self.stored_agents[name] = chain
        return chain


    #
    # def create_agent(self,llm,tools,prompt,name):
    #     if(tools==None):
    #         tools=[]
    #
    #     aux_llm=self.get_one_llm(llm)
    #
    #     template =f"""{prompt}
    #
    #     Você tem acesso às seguintes ferramentas:
    #     {{tools}}
    #
    #     As ferramentas disponíveis são: {{tool_names}}
    #
    #     Quando usar ferramentas, siga este formato de raciocínio:
    #     {{agent_scratchpad}}
    #
    #     Usuário: {{input}}
    #     """
    #
    #     final_prompt = ChatPromptTemplate.from_template(template)
    #
    #     aux_agent=create_react_agent(llm=aux_llm,tools=tools,prompt=final_prompt)
    #     aux_agent.name=name
    #     self.stored_agents[name]=aux_agent
    #     return aux_agent






"""  
llm

agent

tool

callable
"""
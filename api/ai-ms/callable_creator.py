from config import State
from agents_llm_core import AgentLLMCore
from graph_core import GraphCore

class CallableCreator:

    def create_callable_agent(self, agent):
        def node_fn(state):
            # 1) Decide qual texto usar como entrada:
            #    se já tem output de node anterior, usa ele;
            #    senão, usa o input original.
            user_text = state.get("output") or state.get("input", "")

            # 2) Chama o "agente" (no seu caso, um chain simples)
            result = agent.invoke({"input": user_text})

            # 3) Extrai texto da resposta
            if hasattr(result, "content"):          # ChatGroq / ChatOpenAI
                output_text = result.content
            elif isinstance(result, dict) and "output" in result:
                output_text = result["output"]
            else:
                output_text = str(result)

            # 4) Atualiza APENAS o campo output
            return {"output": output_text}

        return node_fn


    # def create_callable_agent(self, agent):
    #     def node_fn(state):
    #         # Importante: passar "intermediate_steps": [] para o ReAct agent
    #         result = agent.invoke({
    #             "input": state["input"],
    #             "intermediate_steps": []
    #         })
    #
    #         # Em muitas versões o resultado vem como dict {"output": "...", ...}
    #         if isinstance(result, dict):
    #             return {"output": result.get("output", str(result))}
    #
    #         # Se vier como string ou outro tipo
    #         return {"output": str(result)}
    #
    #     return node_fn



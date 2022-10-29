import asyncio

done, pending = await asyncio. gather(tasks, timeout=10)